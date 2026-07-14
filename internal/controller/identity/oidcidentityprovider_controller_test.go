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

package identity

import (
	"context"

	"github.com/kubehippie/keycloak-operator/api/common"
	identityv1alpha1 "github.com/kubehippie/keycloak-operator/api/identity/v1alpha1"
	controller "github.com/kubehippie/keycloak-operator/internal/controller"
	"github.com/kubehippie/keycloak-operator/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const testOIDCIdentityProviderAlias = "sample-oidc-idp"
const testOIDCIdentityProviderAuthURL = "https://idp.example.com/oauth2/authorize"
const testOIDCIdentityProviderTokenURL = "https://idp.example.com/oauth2/token"
const testOIDCIdentityProviderClientID = "sample-client-id"

var _ = Describe("OIDCIdentityProvider Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: utils.StandardTestNamespace,
		}
		oidcidentityprovider := &identityv1alpha1.OIDCIdentityProvider{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind OIDCIdentityProvider")
			err := k8sClient.Get(ctx, typeNamespacedName, oidcidentityprovider)
			if err != nil && errors.IsNotFound(err) {
				resource := &identityv1alpha1.OIDCIdentityProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: utils.StandardTestNamespace,
					},
					Spec: identityv1alpha1.OIDCIdentityProviderSpec{
						RealmRef:         &common.RealmRef{Name: "missing-realm"},
						Alias:            testOIDCIdentityProviderAlias,
						AuthorizationURL: testOIDCIdentityProviderAuthURL,
						TokenURL:         testOIDCIdentityProviderTokenURL,
						ClientID:         testOIDCIdentityProviderClientID,
						ClientSecret:     &common.SecretKeyRefOrVal{Value: "sample-secret"},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &identityv1alpha1.OIDCIdentityProvider{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance OIDCIdentityProvider")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should requeue when the Keycloak realm cannot yet be resolved", func() {
			By("Reconciling the created resource")
			controllerReconciler := &OIDCIdentityProviderReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(Equal(controller.FailedKeycloakConnectionRetryPeriod))
		})
	})
})

var _ = Describe("oidcIdentityProviderToGocloak", func() {
	It("maps the identity provider payload", func() {
		provider := &identityv1alpha1.OIDCIdentityProvider{
			Spec: identityv1alpha1.OIDCIdentityProviderSpec{
				Alias:            testOIDCIdentityProviderAlias,
				AuthorizationURL: testOIDCIdentityProviderAuthURL,
				TokenURL:         testOIDCIdentityProviderTokenURL,
				ClientID:         testOIDCIdentityProviderClientID,
				ClientSecret:     &common.SecretKeyRefOrVal{Value: "sample-secret"},
			},
		}

		got, err := oidcIdentityProviderToGocloak(context.Background(), k8sClient, provider, utils.StandardTestNamespace)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Alias).NotTo(BeNil())
		Expect(*got.Alias).To(Equal(testOIDCIdentityProviderAlias))
		Expect(got.Config).To(HaveKeyWithValue("authorizationUrl", testOIDCIdentityProviderAuthURL))
		Expect(got.Config).To(HaveKeyWithValue("tokenUrl", testOIDCIdentityProviderTokenURL))
		Expect(got.Config).To(HaveKeyWithValue("clientId", testOIDCIdentityProviderClientID))
		Expect(got.Config).To(HaveKeyWithValue("clientSecret", "sample-secret"))
	})
})
