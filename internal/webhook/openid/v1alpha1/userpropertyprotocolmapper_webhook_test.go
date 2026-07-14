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
	openidv1alpha1 "github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserPropertyProtocolMapper Webhook", func() {
	var (
		obj       *openidv1alpha1.UserPropertyProtocolMapper
		oldObj    *openidv1alpha1.UserPropertyProtocolMapper
		validator UserPropertyProtocolMapperCustomValidator
		defaulter UserPropertyProtocolMapperCustomDefaulter
	)

	validSpec := func() openidv1alpha1.UserPropertyProtocolMapperSpec {
		return openidv1alpha1.UserPropertyProtocolMapperSpec{
			ClientRef:    &common.ClientRef{Name: testSampleOpenIDClientName},
			Name:         testScopeEmail,
			UserProperty: testScopeEmail,
			ClaimName:    testScopeEmail,
		}
	}

	BeforeEach(func() {
		obj = &openidv1alpha1.UserPropertyProtocolMapper{}
		oldObj = &openidv1alpha1.UserPropertyProtocolMapper{}
		validator = UserPropertyProtocolMapperCustomValidator{}
		defaulter = UserPropertyProtocolMapperCustomDefaulter{}
	})

	It("defaults claimValueType and token booleans", func() {
		obj.Spec = validSpec()
		Expect(defaulter.Default(ctx, obj)).To(Succeed())
		Expect(*obj.Spec.ClaimValueType).To(Equal("String"))
		Expect(*obj.Spec.AddToIDToken).To(BeTrue())
		Expect(*obj.Spec.AddToAccessToken).To(BeTrue())
		Expect(*obj.Spec.AddToUserInfo).To(BeTrue())
	})

	It("admits a valid create", func() {
		obj.Spec = validSpec()
		_, err := validator.ValidateCreate(ctx, obj)
		Expect(err).NotTo(HaveOccurred())
	})

	It("rejects missing userProperty", func() {
		obj.Spec = validSpec()
		obj.Spec.UserProperty = ""
		_, err := validator.ValidateCreate(ctx, obj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("spec.userProperty"))
	})

	It("rejects missing claimName", func() {
		obj.Spec = validSpec()
		obj.Spec.ClaimName = ""
		_, err := validator.ValidateCreate(ctx, obj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("spec.claimName"))
	})

	It("rejects invalid claimValueType", func() {
		obj.Spec = validSpec()
		obj.Spec.ClaimValueType = stringPtr("number")
		_, err := validator.ValidateCreate(ctx, obj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("spec.claimValueType"))
	})

	It("admits updates when clientRef is unchanged", func() {
		oldObj.Spec = validSpec()
		obj.Spec = validSpec()
		_, err := validator.ValidateUpdate(ctx, oldObj, obj)
		Expect(err).NotTo(HaveOccurred())
	})

	It("rejects clientRef changes on update", func() {
		oldObj.Spec = validSpec()
		obj.Spec = validSpec()
		obj.Spec.ClientRef = &common.ClientRef{Name: testOtherClientName}
		_, err := validator.ValidateUpdate(ctx, oldObj, obj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("immutable"))
	})
})
