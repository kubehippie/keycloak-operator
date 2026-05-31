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
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Keycloak Webhook", func() {
	var (
		obj       *v1alpha1.Keycloak
		oldObj    *v1alpha1.Keycloak
		validator KeycloakCustomValidator
		defaulter KeycloakCustomDefaulter
	)

	BeforeEach(func() {
		obj = &v1alpha1.Keycloak{}
		oldObj = &v1alpha1.Keycloak{}
		validator = KeycloakCustomValidator{}
		defaulter = KeycloakCustomDefaulter{}
		Expect(validator).NotTo(BeNil())
		Expect(defaulter).NotTo(BeNil())
	})

	validSpec := func() v1alpha1.KeycloakSpec {
		return v1alpha1.KeycloakSpec{
			URL:       "https://keycloak.example.com",
			RealmName: "master",
			Username:  &common.SecretKeyRefOrVal{Value: "admin"},
			Password:  &common.SecretKeyRefOrVal{Value: "secret"},
		}
	}

	Context("When creating Keycloak under Defaulting Webhook", func() {
		It("Should apply defaults without error", func() {
			obj.Spec = validSpec()
			Expect(defaulter.Default(ctx, obj)).To(Succeed())
		})
	})

	Context("When creating or updating Keycloak under Validating Webhook", func() {
		It("Should admit creation when all required fields are present with inline values", func() {
			obj.Spec = validSpec()
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should admit creation when credentials are supplied as secret references", func() {
			obj.Spec = v1alpha1.KeycloakSpec{
				URL:       "https://keycloak.example.com",
				RealmName: "master",
				Username: &common.SecretKeyRefOrVal{
					SecretKeyRef: &common.SecretKeySelector{
						Name: "kc-secret",
						Key:  "username",
					},
				},
				Password: &common.SecretKeyRefOrVal{
					SecretKeyRef: &common.SecretKeySelector{
						Name: "kc-secret",
						Key:  "password",
					},
				},
			}
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should deny creation when url is missing", func() {
			obj.Spec = validSpec()
			obj.Spec.URL = ""
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.url"))
		})

		It("Should deny creation when url is not a valid URL", func() {
			obj.Spec = validSpec()
			obj.Spec.URL = "not-a-url"
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.url"))
		})

		It("Should deny creation when realmName is missing", func() {
			obj.Spec = validSpec()
			obj.Spec.RealmName = ""
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.realmName"))
		})

		It("Should deny creation when username is nil", func() {
			obj.Spec = validSpec()
			obj.Spec.Username = nil
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.username"))
		})

		It("Should deny creation when username has neither value nor secret ref", func() {
			obj.Spec = validSpec()
			obj.Spec.Username = &common.SecretKeyRefOrVal{}
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.username"))
		})

		It("Should deny creation when password is nil", func() {
			obj.Spec = validSpec()
			obj.Spec.Password = nil
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.password"))
		})

		It("Should deny creation when password has neither value nor secret ref", func() {
			obj.Spec = validSpec()
			obj.Spec.Password = &common.SecretKeyRefOrVal{}
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.password"))
		})

		It("Should admit update when spec is valid", func() {
			obj.Spec = validSpec()
			oldObj.Spec = validSpec()
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should deny update when url becomes invalid", func() {
			oldObj.Spec = validSpec()
			obj.Spec = validSpec()
			obj.Spec.URL = "not-a-url"
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).To(HaveOccurred())
		})

		It("Should admit deletion", func() {
			obj.Spec = validSpec()
			_, err := validator.ValidateDelete(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
