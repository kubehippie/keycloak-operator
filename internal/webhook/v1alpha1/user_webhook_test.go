/*
Copyright 2026 Thomas Boerger <thomas@webhippie.de>.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kubehippie/keycloak-operator/api/common"
	keycloakoperatorwebhippiedev1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
)

var _ = Describe("User Webhook", func() {
	var (
		obj       *keycloakoperatorwebhippiedev1alpha1.User
		oldObj    *keycloakoperatorwebhippiedev1alpha1.User
		validator UserCustomValidator
		defaulter UserCustomDefaulter
	)

	BeforeEach(func() {
		obj = &keycloakoperatorwebhippiedev1alpha1.User{}
		oldObj = &keycloakoperatorwebhippiedev1alpha1.User{}
		validator = UserCustomValidator{}
		defaulter = UserCustomDefaulter{}
		Expect(validator).NotTo(BeNil())
		Expect(defaulter).NotTo(BeNil())
	})

	validSpec := func() keycloakoperatorwebhippiedev1alpha1.UserSpec {
		return keycloakoperatorwebhippiedev1alpha1.UserSpec{
			RealmRef: &common.RealmRef{Kind: "Realm", Name: "my-realm"},
			Username: "johndoe",
		}
	}

	Context("When creating User under Defaulting Webhook", func() {
		It("Should apply defaults without error", func() {
			obj.Spec = validSpec()
			Expect(defaulter.Default(ctx, obj)).To(Succeed())
		})
	})

	Context("When creating User under Validating Webhook", func() {
		It("Should admit creation when all required fields are present", func() {
			obj.Spec = validSpec()
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should deny creation when realmRef is nil", func() {
			obj.Spec = validSpec()
			obj.Spec.RealmRef = nil
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.realmRef.name"))
		})

		It("Should deny creation when realmRef.name is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.RealmRef = &common.RealmRef{Kind: "Realm", Name: ""}
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.realmRef.name"))
		})

		It("Should deny creation when username is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.Username = ""
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.username"))
		})

		It("Should deny creation when username is whitespace only", func() {
			obj.Spec = validSpec()
			obj.Spec.Username = "   "
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.username"))
		})

		It("Should admit deletion", func() {
			obj.Spec = validSpec()
			_, err := validator.ValidateDelete(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When updating User under Validating Webhook", func() {
		It("Should admit update when username is unchanged", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should deny update when username is changed", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			obj.Spec.Username = "janedoe"
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("immutable"))
		})

		It("Should deny update when realmRef is removed", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			obj.Spec.RealmRef = nil
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.realmRef.name"))
		})
	})
})
