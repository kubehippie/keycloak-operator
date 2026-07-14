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
	"github.com/kubehippie/keycloak-operator/api/common"
	"github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientScope Webhook", func() {
	var (
		obj       *v1alpha1.ClientScope
		oldObj    *v1alpha1.ClientScope
		validator ClientScopeCustomValidator
		defaulter ClientScopeCustomDefaulter
	)

	BeforeEach(func() {
		obj = &v1alpha1.ClientScope{}
		oldObj = &v1alpha1.ClientScope{}
		validator = ClientScopeCustomValidator{}
		defaulter = ClientScopeCustomDefaulter{}
		Expect(validator).NotTo(BeNil())
		Expect(defaulter).NotTo(BeNil())
	})

	validSpec := func() v1alpha1.ClientScopeSpec {
		return v1alpha1.ClientScopeSpec{
			RealmRef: &common.RealmRef{Name: testSampleRealmName},
			Name:     "profile-extra",
		}
	}

	Context("When creating ClientScope under Defaulting Webhook", func() {
		It("Should default includeInTokenScope to true", func() {
			obj.Spec = validSpec()
			Expect(defaulter.Default(ctx, obj)).To(Succeed())
			Expect(obj.Spec.IncludeInTokenScope).NotTo(BeNil())
			Expect(*obj.Spec.IncludeInTokenScope).To(BeTrue())
		})
	})

	Context("When creating ClientScope under Validating Webhook", func() {
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
			Expect(err.Error()).To(ContainSubstring("spec.realmRef"))
		})

		It("Should deny creation when realmRef.name is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.RealmRef = &common.RealmRef{Name: ""}
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.realmRef.name"))
		})

		It("Should deny creation when name is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.Name = ""
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.name"))
		})

		It("Should admit deletion", func() {
			obj.Spec = validSpec()
			_, err := validator.ValidateDelete(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When updating ClientScope under Validating Webhook", func() {
		It("Should admit update when realmRef is unchanged", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should deny update when realmRef is changed", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			obj.Spec.RealmRef = &common.RealmRef{Name: "other-realm"}
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("immutable"))
		})
	})
})
