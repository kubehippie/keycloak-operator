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
	identityv1alpha1 "github.com/kubehippie/keycloak-operator/api/identity/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AttributeImporterMapper Webhook", func() {
	var (
		obj       *identityv1alpha1.AttributeImporterMapper
		oldObj    *identityv1alpha1.AttributeImporterMapper
		validator AttributeImporterMapperCustomValidator
		defaulter AttributeImporterMapperCustomDefaulter
	)

	BeforeEach(func() {
		obj = &identityv1alpha1.AttributeImporterMapper{}
		oldObj = &identityv1alpha1.AttributeImporterMapper{}
		validator = AttributeImporterMapperCustomValidator{}
		defaulter = AttributeImporterMapperCustomDefaulter{}
		Expect(validator).NotTo(BeNil())
		Expect(defaulter).NotTo(BeNil())
	})

	validSpec := func() identityv1alpha1.AttributeImporterMapperSpec {
		return identityv1alpha1.AttributeImporterMapperSpec{
			IdentityProviderRef: &common.IdentityProviderRef{Name: testSampleIdentityProviderName},
			Name:                "import-preferred-username",
			UserAttribute:       "username",
			ClaimName:           "preferred_username",
		}
	}

	Context("When creating AttributeImporterMapper under Defaulting Webhook", func() {
		It("Should apply defaults without error", func() {
			obj.Spec = validSpec()
			Expect(defaulter.Default(ctx, obj)).To(Succeed())
		})
	})

	Context("When creating AttributeImporterMapper under Validating Webhook", func() {
		It("Should admit creation when all required fields are present", func() {
			obj.Spec = validSpec()
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should deny creation when identityProviderRef is nil", func() {
			obj.Spec = validSpec()
			obj.Spec.IdentityProviderRef = nil
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.identityProviderRef"))
		})

		It("Should deny creation when identityProviderRef.name is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.IdentityProviderRef = &common.IdentityProviderRef{Name: ""}
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.identityProviderRef.name"))
		})

		It("Should deny creation when name is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.Name = ""
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.name"))
		})

		It("Should deny creation when userAttribute is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.UserAttribute = ""
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.userAttribute"))
		})

		It("Should deny creation when claimName is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.ClaimName = ""
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.claimName"))
		})

		It("Should admit deletion", func() {
			obj.Spec = validSpec()
			_, err := validator.ValidateDelete(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When updating AttributeImporterMapper under Validating Webhook", func() {
		It("Should admit update when identityProviderRef is unchanged", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should deny update when identityProviderRef is changed", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			obj.Spec.IdentityProviderRef = &common.IdentityProviderRef{Name: testOtherIdentityProviderName}
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("immutable"))
		})

		It("Should deny update when identityProviderRef is removed", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			obj.Spec.IdentityProviderRef = nil
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.identityProviderRef"))
		})
	})
})
