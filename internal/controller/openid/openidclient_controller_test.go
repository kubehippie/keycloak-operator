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
	controller "github.com/kubehippie/keycloak-operator/internal/controller"
	"github.com/kubehippie/keycloak-operator/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const testOpenIDClientClientID = "sample-client"

var _ = Describe("OpenIDClient Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: utils.StandardTestNamespace,
		}
		openidclient := &v1alpha1.OpenIDClient{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind OpenIDClient")
			err := k8sClient.Get(ctx, typeNamespacedName, openidclient)
			if err != nil && errors.IsNotFound(err) {
				resource := &v1alpha1.OpenIDClient{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: utils.StandardTestNamespace,
					},
					Spec: v1alpha1.OpenIDClientSpec{
						RealmRef:   &common.RealmRef{Name: "missing-realm"},
						ClientID:   testOpenIDClientClientID,
						AccessType: "CONFIDENTIAL",
						ClientSecret: &common.SecretKeyRefOrVal{
							Value: "sample-secret",
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &v1alpha1.OpenIDClient{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance OpenIDClient")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should requeue when the Keycloak realm cannot yet be resolved", func() {
			By("Reconciling the created resource")
			controllerReconciler := &OpenIDClientReconciler{
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

var _ = Describe("openIDClientToGocloak", func() {
	It("maps the client payload", func() {
		client := &v1alpha1.OpenIDClient{
			Spec: v1alpha1.OpenIDClientSpec{
				ClientID:   testOpenIDClientClientID,
				AccessType: "CONFIDENTIAL",
				ClientSecret: &common.SecretKeyRefOrVal{
					Value: "sample-secret",
				},
			},
		}

		got, err := openIDClientToGocloak(context.Background(), k8sClient, client, utils.StandardTestNamespace)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.ClientID).NotTo(BeNil())
		Expect(*got.ClientID).To(Equal(testOpenIDClientClientID))
		Expect(got.Secret).NotTo(BeNil())
		Expect(*got.Secret).To(Equal("sample-secret"))
		Expect(got.PublicClient).NotTo(BeNil())
		Expect(*got.PublicClient).To(BeFalse())
		Expect(got.BearerOnly).NotTo(BeNil())
		Expect(*got.BearerOnly).To(BeFalse())
	})

	It("maps a public client without a secret", func() {
		client := &v1alpha1.OpenIDClient{
			Spec: v1alpha1.OpenIDClientSpec{
				ClientID:   testOpenIDClientClientID,
				AccessType: "PUBLIC",
			},
		}

		got, err := openIDClientToGocloak(context.Background(), k8sClient, client, utils.StandardTestNamespace)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.PublicClient).NotTo(BeNil())
		Expect(*got.PublicClient).To(BeTrue())
		Expect(got.Secret).To(BeNil())
	})
})
