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

package openid

import (
	"context"

	"github.com/kubehippie/keycloak-operator/api/common"
	"github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	"github.com/kubehippie/keycloak-operator/internal/controller"
	"github.com/kubehippie/keycloak-operator/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const testClientScopeName = "profile-extra"

var _ = Describe("ClientScope Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: utils.StandardTestNamespace,
		}
		clientscope := &v1alpha1.ClientScope{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind ClientScope")
			err := k8sClient.Get(ctx, typeNamespacedName, clientscope)
			if err != nil && errors.IsNotFound(err) {
				resource := &v1alpha1.ClientScope{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: utils.StandardTestNamespace,
					},
					Spec: v1alpha1.ClientScopeSpec{
						RealmRef: &common.RealmRef{Name: "missing-realm"},
						Name:     testClientScopeName,
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &v1alpha1.ClientScope{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance ClientScope")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should requeue when the realm cannot yet be resolved", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ClientScopeReconciler{
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

var _ = Describe("clientScopeToGocloak", func() {
	It("maps the client scope payload with defaults", func() {
		includeInTokenScope := true
		scope := &v1alpha1.ClientScope{
			Spec: v1alpha1.ClientScopeSpec{
				Name:                testClientScopeName,
				IncludeInTokenScope: &includeInTokenScope,
			},
		}

		got := clientScopeToGocloak(scope)
		Expect(got.Name).NotTo(BeNil())
		Expect(*got.Name).To(Equal(testClientScopeName))
		Expect(got.Protocol).NotTo(BeNil())
		Expect(*got.Protocol).To(Equal("openid-connect"))
		Expect(got.ClientScopeAttributes.IncludeInTokenScope).NotTo(BeNil())
		Expect(*got.ClientScopeAttributes.IncludeInTokenScope).To(Equal("true"))
		Expect(got.ClientScopeAttributes.DisplayOnConsentScreen).NotTo(BeNil())
		Expect(*got.ClientScopeAttributes.DisplayOnConsentScreen).To(Equal("false"))
	})

	It("enables the consent screen when consentScreenText is set", func() {
		consentText := "Grants access to your profile"
		scope := &v1alpha1.ClientScope{
			Spec: v1alpha1.ClientScopeSpec{
				Name:              testClientScopeName,
				ConsentScreenText: &consentText,
			},
		}

		got := clientScopeToGocloak(scope)
		Expect(got.ClientScopeAttributes.ConsentScreenText).To(Equal(&consentText))
		Expect(got.ClientScopeAttributes.DisplayOnConsentScreen).NotTo(BeNil())
		Expect(*got.ClientScopeAttributes.DisplayOnConsentScreen).To(Equal("true"))
	})
})
