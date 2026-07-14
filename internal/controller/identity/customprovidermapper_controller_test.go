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

const testCustomProviderMapperName = "custom-username-mapper"
const testCustomProviderMapperType = "oidc-username-idp-mapper"

var _ = Describe("CustomProviderMapper Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: utils.StandardTestNamespace,
		}
		customprovidermapper := &identityv1alpha1.CustomProviderMapper{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind CustomProviderMapper")
			err := k8sClient.Get(ctx, typeNamespacedName, customprovidermapper)
			if err != nil && errors.IsNotFound(err) {
				resource := &identityv1alpha1.CustomProviderMapper{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: utils.StandardTestNamespace,
					},
					Spec: identityv1alpha1.CustomProviderMapperSpec{
						IdentityProviderRef: &common.IdentityProviderRef{Name: testMissingIdentityProviderName},
						Name:                testCustomProviderMapperName,
						MapperType:          testCustomProviderMapperType,
						Config: map[string]string{
							"template": "${CLAIM.preferred_username}",
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &identityv1alpha1.CustomProviderMapper{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance CustomProviderMapper")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should requeue when the identity provider cannot yet be resolved", func() {
			By("Reconciling the created resource")
			controllerReconciler := &CustomProviderMapperReconciler{
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

var _ = Describe("customProviderMapperToGocloak", func() {
	It("maps the identity provider mapper payload", func() {
		mapper := &identityv1alpha1.CustomProviderMapper{
			Spec: identityv1alpha1.CustomProviderMapperSpec{
				Name:       testCustomProviderMapperName,
				MapperType: testCustomProviderMapperType,
				Config: map[string]string{
					"template": "${CLAIM.preferred_username}",
				},
			},
		}

		got := customProviderMapperToGocloak(mapper, testIdentityProviderAlias)
		Expect(got.Name).NotTo(BeNil())
		Expect(*got.Name).To(Equal(testCustomProviderMapperName))
		Expect(got.IdentityProviderMapper).NotTo(BeNil())
		Expect(*got.IdentityProviderMapper).To(Equal(testCustomProviderMapperType))
		Expect(got.IdentityProviderAlias).NotTo(BeNil())
		Expect(*got.IdentityProviderAlias).To(Equal(testIdentityProviderAlias))
		Expect(got.Config).To(HaveKeyWithValue("template", "${CLAIM.preferred_username}"))
	})
})
