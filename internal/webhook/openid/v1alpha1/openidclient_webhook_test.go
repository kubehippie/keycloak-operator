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

var _ = Describe("OpenIDClient Webhook", func() {
	const testOpenIDClientClientID = "sample-client"

	var (
		obj       *openidv1alpha1.OpenIDClient
		oldObj    *openidv1alpha1.OpenIDClient
		validator OpenIDClientCustomValidator
		defaulter OpenIDClientCustomDefaulter
	)

	validSpec := func() openidv1alpha1.OpenIDClientSpec {
		return openidv1alpha1.OpenIDClientSpec{
			RealmRef:   &common.RealmRef{Name: testSampleRealmName},
			ClientID:   testOpenIDClientClientID,
			AccessType: accessTypeConfidential,
			ClientSecret: &common.SecretKeyRefOrVal{
				Value: "sample-secret",
			},
		}
	}

	BeforeEach(func() {
		obj = &openidv1alpha1.OpenIDClient{}
		oldObj = &openidv1alpha1.OpenIDClient{}
		validator = OpenIDClientCustomValidator{}
		defaulter = OpenIDClientCustomDefaulter{}
	})

	Context("When creating OpenIDClient under Defaulting Webhook", func() {
		It("Should apply defaults", func() {
			obj.Spec = openidv1alpha1.OpenIDClientSpec{
				RealmRef: &common.RealmRef{Name: testSampleRealmName},
				ClientID: testOpenIDClientClientID,
			}
			Expect(defaulter.Default(ctx, obj)).To(Succeed())
			Expect(obj.Spec.AccessType).To(Equal(accessTypeConfidential))
			Expect(obj.Spec.Enabled).NotTo(BeNil())
			Expect(*obj.Spec.Enabled).To(BeTrue())
			Expect(obj.Spec.StandardFlowEnabled).NotTo(BeNil())
			Expect(*obj.Spec.StandardFlowEnabled).To(BeTrue())
			Expect(obj.Spec.ClientAuthenticatorType).NotTo(BeNil())
			Expect(*obj.Spec.ClientAuthenticatorType).To(Equal("client-secret"))
		})
	})

	Context("When creating or updating OpenIDClient under Validating Webhook", func() {
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

		It("Should deny creation when clientID is empty", func() {
			obj.Spec = validSpec()
			obj.Spec.ClientID = ""
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.clientID"))
		})

		It("Should deny creation when accessType is invalid", func() {
			obj.Spec = validSpec()
			obj.Spec.AccessType = "INVALID"
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.accessType"))
		})

		It("Should deny creation when accessType is CONFIDENTIAL without a clientSecret", func() {
			obj.Spec = validSpec()
			obj.Spec.ClientSecret = nil
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.clientSecret"))
		})

		It("Should admit deletion", func() {
			obj.Spec = validSpec()
			_, err := validator.ValidateDelete(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("When updating OpenIDClient under Validating Webhook", func() {
		It("Should admit update when clientID is unchanged", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should deny update when clientID is changed", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			obj.Spec.ClientID = "other-client-id"
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("immutable"))
		})
	})
})
