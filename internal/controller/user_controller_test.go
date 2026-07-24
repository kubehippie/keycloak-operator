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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const testUserUsername = "alice"

var _ = Describe("User Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: utils.StandardTestNamespace,
		}
		user := &v1alpha1.User{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind User")
			err := k8sClient.Get(ctx, typeNamespacedName, user)
			if err != nil && errors.IsNotFound(err) {
				resource := &v1alpha1.User{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: utils.StandardTestNamespace,
					},
					Spec: v1alpha1.UserSpec{
						RealmRef: &common.RealmRef{
							Name: utils.StandardTestRealmName,
						},
						Username: testUserUsername,
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &v1alpha1.User{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance User")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &UserReconciler{
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

var _ = Describe("userToGocloak", func() {
	It("maps username", func() {
		u := &v1alpha1.User{
			Spec: v1alpha1.UserSpec{Username: testUserUsername},
		}
		g := userToGocloak(u)
		Expect(g.Username).NotTo(BeNil())
		Expect(*g.Username).To(Equal(testUserUsername))
	})

	It("maps optional fields when set", func() {
		email := "alice@example.com"
		enabled := true
		verified := false
		first := "Alice"
		last := "Smith"

		u := &v1alpha1.User{
			Spec: v1alpha1.UserSpec{
				Username:      testUserUsername,
				Email:         &email,
				Enabled:       &enabled,
				EmailVerified: &verified,
				FirstName:     &first,
				LastName:      &last,
			},
		}
		g := userToGocloak(u)
		Expect(*g.Email).To(Equal(email))
		Expect(*g.Enabled).To(BeTrue())
		Expect(*g.EmailVerified).To(BeFalse())
		Expect(*g.FirstName).To(Equal(first))
		Expect(*g.LastName).To(Equal(last))
	})

	It("maps attributes when set", func() {
		u := &v1alpha1.User{
			Spec: v1alpha1.UserSpec{
				Username:   testUserUsername,
				Attributes: map[string][]string{"department": {"engineering"}},
			},
		}
		g := userToGocloak(u)
		Expect(g.Attributes).NotTo(BeNil())
		Expect(g.Attributes["department"]).To(ConsistOf("engineering"))
	})

	It("maps requiredActions when set", func() {
		u := &v1alpha1.User{
			Spec: v1alpha1.UserSpec{
				Username:        testUserUsername,
				RequiredActions: []string{"UPDATE_PASSWORD"},
			},
		}
		g := userToGocloak(u)
		Expect(g.RequiredActions).NotTo(BeNil())
		Expect(g.RequiredActions).To(ConsistOf("UPDATE_PASSWORD"))
	})

	It("leaves optional fields nil when not set", func() {
		u := &v1alpha1.User{
			Spec: v1alpha1.UserSpec{Username: "bob"},
		}
		g := userToGocloak(u)
		Expect(g.Email).To(BeNil())
		Expect(g.Enabled).To(BeNil())
		Expect(g.EmailVerified).To(BeNil())
		Expect(g.FirstName).To(BeNil())
		Expect(g.LastName).To(BeNil())
		Expect(g.Attributes).To(BeNil())
		Expect(g.RequiredActions).To(BeNil())
	})
})

var _ = Describe("diffNames", func() {
	It("returns all desired names as toAdd when nothing was previously applied", func() {
		toAdd, toRemove := diffNames([]string{"a", "b"}, nil)
		Expect(toAdd).To(ConsistOf("a", "b"))
		Expect(toRemove).To(BeEmpty())
	})

	It("returns previously applied names as toRemove when nothing is desired anymore", func() {
		toAdd, toRemove := diffNames(nil, []string{"a", "b"})
		Expect(toAdd).To(BeEmpty())
		Expect(toRemove).To(ConsistOf("a", "b"))
	})

	It("returns only the difference when desired and previous overlap", func() {
		toAdd, toRemove := diffNames([]string{"a", "c"}, []string{"a", "b"})
		Expect(toAdd).To(ConsistOf("c"))
		Expect(toRemove).To(ConsistOf("b"))
	})

	It("returns no changes when desired matches previously applied", func() {
		toAdd, toRemove := diffNames([]string{"a", "b"}, []string{"a", "b"})
		Expect(toAdd).To(BeEmpty())
		Expect(toRemove).To(BeEmpty())
	})
})

var _ = Describe("UserReconciler.setPassword", func() {
	ctx := context.Background()

	It("is a no-op when spec.password is not set", func() {
		reconciler := &UserReconciler{Client: k8sClient}
		u := &v1alpha1.User{Spec: v1alpha1.UserSpec{Username: testUserUsername}}
		err := reconciler.setPassword(ctx, u, nil, "some-keycloak-id")
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns an error when spec.password cannot be resolved", func() {
		reconciler := &UserReconciler{Client: k8sClient}
		u := &v1alpha1.User{
			ObjectMeta: metav1.ObjectMeta{Namespace: utils.StandardTestNamespace},
			Spec: v1alpha1.UserSpec{
				Username: testUserUsername,
				Password: &common.SecretKeyRefOrVal{
					SecretKeyRef: &common.SecretKeySelector{Name: "no-such-secret", Key: "password"},
				},
			},
		}
		err := reconciler.setPassword(ctx, u, nil, "some-keycloak-id")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("unable to resolve password"))
	})
})
