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

package controller

import (
	"context"

	"github.com/kubehippie/keycloak-operator/api/common"
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	"github.com/kubehippie/keycloak-operator/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("keycloak_client helpers", func() {
	const testNS = "default"

	ctx := context.Background()

	Describe("resolveSecretKeyRefOrVal", func() {
		It("returns inline value directly without any Secret lookup", func() {
			val, err := resolveSecretKeyRefOrVal(ctx, k8sClient, &common.SecretKeyRefOrVal{Value: "hello"}, testNS)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal("hello"))
		})

		It("errors when ref is nil", func() {
			_, err := resolveSecretKeyRefOrVal(ctx, k8sClient, nil, testNS)
			Expect(err).To(HaveOccurred())
		})

		It("errors when neither value nor secretKeyRef is set", func() {
			_, err := resolveSecretKeyRefOrVal(ctx, k8sClient, &common.SecretKeyRefOrVal{}, testNS)
			Expect(err).To(HaveOccurred())
		})

		It("errors when secretKeyRef.name is empty", func() {
			_, err := resolveSecretKeyRefOrVal(ctx, k8sClient, &common.SecretKeyRefOrVal{
				SecretKeyRef: &common.SecretKeySelector{Key: "just-key-without-name"},
			}, testNS)
			Expect(err).To(HaveOccurred())
		})

		It("errors when secretKeyRef.key is empty", func() {
			_, err := resolveSecretKeyRefOrVal(ctx, k8sClient, &common.SecretKeyRefOrVal{
				SecretKeyRef: &common.SecretKeySelector{Name: "just-name-without-key"},
			}, testNS)
			Expect(err).To(HaveOccurred())
		})

		Context("with a backing Secret", func() {
			const secretName = "kclient-test-secret"
			const keyName = "mykey"

			BeforeEach(func() {
				secret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: testNS},
					Data:       map[string][]byte{keyName: []byte("myvalue")},
				}
				err := k8sClient.Create(ctx, secret)
				if err != nil {
					By("secret already exists, skipping creation")
				}
			})

			AfterEach(func() {
				secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: testNS}}
				_ = k8sClient.Delete(ctx, secret)
			})

			It("reads value from Secret using namespace fallback when ref.Namespace is empty", func() {
				val, err := resolveSecretKeyRefOrVal(ctx, k8sClient, &common.SecretKeyRefOrVal{
					SecretKeyRef: &common.SecretKeySelector{Name: secretName, Key: keyName},
				}, testNS)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal("myvalue"))
			})

			It("reads value from Secret when namespace is given explicitly in the ref", func() {
				val, err := resolveSecretKeyRefOrVal(ctx, k8sClient, &common.SecretKeyRefOrVal{
					SecretKeyRef: &common.SecretKeySelector{Name: secretName, Key: keyName, Namespace: testNS},
				}, "other-ns")
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal("myvalue"))
			})

			It("errors when the key is absent in the Secret", func() {
				_, err := resolveSecretKeyRefOrVal(ctx, k8sClient, &common.SecretKeyRefOrVal{
					SecretKeyRef: &common.SecretKeySelector{Name: secretName, Key: "notexist"},
				}, testNS)
				Expect(err).To(HaveOccurred())
			})

			It("errors when the Secret does not exist in the resolved namespace", func() {
				_, err := resolveSecretKeyRefOrVal(ctx, k8sClient, &common.SecretKeyRefOrVal{
					SecretKeyRef: &common.SecretKeySelector{Name: "no-such-secret", Key: keyName},
				}, testNS)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("keycloakSessionForKeycloak namespace fallback", func() {
		const kcName = "ns-fallback-kc"

		BeforeEach(func() {
			kc := &v1alpha1.Keycloak{
				ObjectMeta: metav1.ObjectMeta{Name: kcName, Namespace: testNS},
				Spec: v1alpha1.KeycloakSpec{
					URL:       "http://127.0.0.1:19999",
					RealmName: utils.StandardTestRealmName,
					Username:  &common.SecretKeyRefOrVal{Value: "root"},
					Password:  &common.SecretKeyRefOrVal{Value: "r4nd0m"},
				},
			}
			err := k8sClient.Create(ctx, kc)
			if err != nil {
				By("keycloak already exists, skipping creation")
			}
		})

		AfterEach(func() {
			kc := &v1alpha1.Keycloak{ObjectMeta: metav1.ObjectMeta{Name: kcName, Namespace: testNS}}
			_ = k8sClient.Delete(ctx, kc)
		})

		It("resolves Keycloak via namespace fallback when ref.Namespace is empty", func() {
			// The Keycloak resource lives in testNS; the ref carries no namespace so it
			// should fall back to defaultNamespace=testNS and find the resource.
			// The session will still fail (no real Keycloak at 127.0.0.1:19999), but
			// the error must NOT be a resource-not-found error.
			_, err := keycloakSessionForKeycloak(ctx, k8sClient,
				&common.KeycloakRef{Name: kcName},
				testNS,
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).NotTo(ContainSubstring("unable to fetch keycloak"))
		})

		It("fails resource lookup when ref.Namespace points to a different namespace", func() {
			_, err := keycloakSessionForKeycloak(ctx, k8sClient,
				&common.KeycloakRef{Name: kcName, Namespace: "other-ns"},
				testNS,
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unable to fetch keycloak"))
		})
	})

	Describe("keycloakSessionForRealm namespace fallback", func() {
		const kcName = "realm-ns-fallback-kc"
		const realmResName = "realm-ns-fallback-realm"

		BeforeEach(func() {
			kc := &v1alpha1.Keycloak{
				ObjectMeta: metav1.ObjectMeta{Name: kcName, Namespace: testNS},
				Spec: v1alpha1.KeycloakSpec{
					URL:       "http://127.0.0.1:19999",
					RealmName: utils.StandardTestRealmName,
					Username:  &common.SecretKeyRefOrVal{Value: "superadmin"},
					Password:  &common.SecretKeyRefOrVal{Value: "5up3r53cr37"},
				},
			}
			if err := k8sClient.Create(ctx, kc); err != nil {
				By("keycloak already exists, skipping creation")
			}

			realm := &v1alpha1.Realm{
				ObjectMeta: metav1.ObjectMeta{Name: realmResName, Namespace: testNS},
				Spec: v1alpha1.RealmSpec{
					// keycloakRef carries no namespace → will fall back to testNS
					KeycloakRef: &common.KeycloakRef{Name: kcName},
					Name:        utils.StandardTestRealmName,
				},
			}
			if err := k8sClient.Create(ctx, realm); err != nil {
				By("realm already exists, skipping creation")
			}
		})

		AfterEach(func() {
			realm := &v1alpha1.Realm{ObjectMeta: metav1.ObjectMeta{Name: realmResName, Namespace: testNS}}
			_ = k8sClient.Delete(ctx, realm)
			kc := &v1alpha1.Keycloak{ObjectMeta: metav1.ObjectMeta{Name: kcName, Namespace: testNS}}
			_ = k8sClient.Delete(ctx, kc)
		})

		It("resolves Realm via namespace fallback when ref.Namespace is empty", func() {
			// Realm lives in testNS; ref carries no namespace → falls back to testNS.
			// Session creation will fail (no live Keycloak), but NOT with a realm-not-found error.
			_, err := keycloakSessionForRealm(ctx, k8sClient,
				&common.RealmRef{Name: realmResName},
				testNS,
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).NotTo(ContainSubstring("unable to fetch realm"))
		})

		It("fails realm lookup when ref.Namespace points to a different namespace", func() {
			_, err := keycloakSessionForRealm(ctx, k8sClient,
				&common.RealmRef{Name: realmResName, Namespace: "other-ns"},
				testNS,
			)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unable to fetch realm"))
		})
	})
})
