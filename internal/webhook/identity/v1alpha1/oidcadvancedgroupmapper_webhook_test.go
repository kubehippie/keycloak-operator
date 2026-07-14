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

var _ = Describe("OIDCAdvancedGroupMapper Webhook", func() {
	var (
		obj       *identityv1alpha1.OIDCAdvancedGroupMapper
		oldObj    *identityv1alpha1.OIDCAdvancedGroupMapper
		validator OIDCAdvancedGroupMapperCustomValidator
		defaulter OIDCAdvancedGroupMapperCustomDefaulter
	)

	BeforeEach(func() {
		obj = &identityv1alpha1.OIDCAdvancedGroupMapper{}
		oldObj = &identityv1alpha1.OIDCAdvancedGroupMapper{}
		validator = OIDCAdvancedGroupMapperCustomValidator{}
		defaulter = OIDCAdvancedGroupMapperCustomDefaulter{}
		Expect(validator).NotTo(BeNil())
		Expect(defaulter).NotTo(BeNil())
	})

	validSpec := func() identityv1alpha1.OIDCAdvancedGroupMapperSpec {
		return identityv1alpha1.OIDCAdvancedGroupMapperSpec{
			IdentityProviderRef: &common.IdentityProviderRef{Name: testSampleIdentityProviderName},
			Name:                "engineering-group-mapper",
			Claims:              map[string]string{"department": "engineering"},
			Group:               "/engineering",
		}
	}

	Context("When creating OIDCAdvancedGroupMapper under Defaulting Webhook", func() {
		It("Should default claimValuesRegex to false", func() {
			obj.Spec = validSpec()
			Expect(defaulter.Default(ctx, obj)).To(Succeed())
			Expect(obj.Spec.ClaimValuesRegex).NotTo(BeNil())
			Expect(*obj.Spec.ClaimValuesRegex).To(BeFalse())
		})
	})

	Context("When creating OIDCAdvancedGroupMapper under Validating Webhook", func() {
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

		It("Should deny creation when claims is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.Claims = nil
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.claims"))
		})

		It("Should deny creation when group is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.Group = ""
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.group"))
		})

		It("Should admit deletion", func() {
			obj.Spec = validSpec()
			_, err := validator.ValidateDelete(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When updating OIDCAdvancedGroupMapper under Validating Webhook", func() {
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
	})
})
