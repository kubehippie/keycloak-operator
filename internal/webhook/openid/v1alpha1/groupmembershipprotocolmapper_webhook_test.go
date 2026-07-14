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

var _ = Describe("GroupMembershipProtocolMapper Webhook", func() {
	var (
		obj       *openidv1alpha1.GroupMembershipProtocolMapper
		oldObj    *openidv1alpha1.GroupMembershipProtocolMapper
		validator GroupMembershipProtocolMapperCustomValidator
		defaulter GroupMembershipProtocolMapperCustomDefaulter
	)

	validSpec := func() openidv1alpha1.GroupMembershipProtocolMapperSpec {
		return openidv1alpha1.GroupMembershipProtocolMapperSpec{
			ClientRef: &common.ClientRef{Name: testSampleOpenIDClientName},
			Name:      "groups",
			ClaimName: "groups",
			FullPath:  nil,
		}
	}

	BeforeEach(func() {
		obj = &openidv1alpha1.GroupMembershipProtocolMapper{}
		oldObj = &openidv1alpha1.GroupMembershipProtocolMapper{}
		validator = GroupMembershipProtocolMapperCustomValidator{}
		defaulter = GroupMembershipProtocolMapperCustomDefaulter{}
	})

	It("defaults all booleans to true", func() {
		obj.Spec = validSpec()
		Expect(defaulter.Default(ctx, obj)).To(Succeed())
		Expect(*obj.Spec.FullPath).To(BeTrue())
		Expect(*obj.Spec.AddToIDToken).To(BeTrue())
		Expect(*obj.Spec.AddToAccessToken).To(BeTrue())
		Expect(*obj.Spec.AddToUserInfo).To(BeTrue())
		Expect(*obj.Spec.AddToTokenIntrospection).To(BeTrue())
	})

	It("admits a valid create", func() {
		obj.Spec = validSpec()
		_, err := validator.ValidateCreate(ctx, obj)
		Expect(err).NotTo(HaveOccurred())
	})

	It("rejects missing clientRef", func() {
		obj.Spec = validSpec()
		obj.Spec.ClientRef = nil
		_, err := validator.ValidateCreate(ctx, obj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("spec.clientRef"))
	})

	It("rejects empty name", func() {
		obj.Spec = validSpec()
		obj.Spec.Name = ""
		_, err := validator.ValidateCreate(ctx, obj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("spec.name"))
	})

	It("rejects empty claimName", func() {
		obj.Spec = validSpec()
		obj.Spec.ClaimName = "  "
		_, err := validator.ValidateCreate(ctx, obj)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("spec.claimName"))
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
