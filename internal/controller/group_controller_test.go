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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubehippie/keycloak-operator/api/common"
	keycloakoperatorwebhippiedev1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
)

var _ = Describe("Group Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		group := &keycloakoperatorwebhippiedev1alpha1.Group{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind Group")
			err := k8sClient.Get(ctx, typeNamespacedName, group)
			if err != nil && errors.IsNotFound(err) {
				resource := &keycloakoperatorwebhippiedev1alpha1.Group{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: keycloakoperatorwebhippiedev1alpha1.GroupSpec{
						RealmRef: &common.RealmRef{
							Name: "test-realm",
						},
						Name: "test-group",
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &keycloakoperatorwebhippiedev1alpha1.Group{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Group")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &GroupReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("groupToGocloak", func() {
	It("maps name", func() {
		g := &keycloakoperatorwebhippiedev1alpha1.Group{
			Spec: keycloakoperatorwebhippiedev1alpha1.GroupSpec{Name: "admins"},
		}
		got := groupToGocloak(g)
		Expect(got.Name).NotTo(BeNil())
		Expect(*got.Name).To(Equal("admins"))
	})

	It("maps attributes when set", func() {
		g := &keycloakoperatorwebhippiedev1alpha1.Group{
			Spec: keycloakoperatorwebhippiedev1alpha1.GroupSpec{
				Name:       "admins",
				Attributes: map[string][]string{"team": {"platform"}},
			},
		}
		got := groupToGocloak(g)
		Expect(got.Attributes).NotTo(BeNil())
		Expect((*got.Attributes)["team"]).To(ConsistOf("platform"))
	})

	It("maps realmRoles when set", func() {
		g := &keycloakoperatorwebhippiedev1alpha1.Group{
			Spec: keycloakoperatorwebhippiedev1alpha1.GroupSpec{
				Name:       "admins",
				RealmRoles: []string{"admin", "viewer"},
			},
		}
		got := groupToGocloak(g)
		Expect(got.RealmRoles).NotTo(BeNil())
		Expect(*got.RealmRoles).To(ConsistOf("admin", "viewer"))
	})

	It("maps clientRoles when set", func() {
		g := &keycloakoperatorwebhippiedev1alpha1.Group{
			Spec: keycloakoperatorwebhippiedev1alpha1.GroupSpec{
				Name:        "admins",
				ClientRoles: map[string][]string{"my-client": {"role-a"}},
			},
		}
		got := groupToGocloak(g)
		Expect(got.ClientRoles).NotTo(BeNil())
		Expect((*got.ClientRoles)["my-client"]).To(ConsistOf("role-a"))
	})

	It("leaves optional fields nil when not set", func() {
		g := &keycloakoperatorwebhippiedev1alpha1.Group{
			Spec: keycloakoperatorwebhippiedev1alpha1.GroupSpec{Name: "empty"},
		}
		got := groupToGocloak(g)
		Expect(got.Attributes).To(BeNil())
		Expect(got.RealmRoles).To(BeNil())
		Expect(got.ClientRoles).To(BeNil())
	})
})
