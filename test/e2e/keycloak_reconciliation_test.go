//go:build e2e
// +build e2e

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

package e2e

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kubehippie/keycloak-operator/test/utils"
)

// keycloakServiceURL is the in-cluster address of the test Keycloak instance.
const keycloakServiceURL = "http://keycloak.keycloak-system.svc.cluster.local"

var _ = Describe("Keycloak Reconciliation", Ordered, func() {
	const (
		testNamespace  = namespace
		keycloakCRName = "e2e-keycloak"
		realmCRName    = "e2e-realm"
		realmName      = "e2e-realm"
		userCRName     = "e2e-user"
		groupCRName    = "e2e-group"

		reconcileTimeout = 3 * time.Minute
		reconcilePolling = 5 * time.Second
	)

	BeforeAll(func() {
		if keycloakSkip {
			Skip("Keycloak reconciliation tests skipped")
		}
	})

	AfterAll(func() {
		if keycloakSkip {
			return
		}

		for _, args := range [][]string{
			{"delete", "users.keycloak-operator.webhippie.de", userCRName, "-n", testNamespace, "--ignore-not-found", "--wait=false"},
			{"delete", "groups.keycloak-operator.webhippie.de", groupCRName, "-n", testNamespace, "--ignore-not-found", "--wait=false"},
			{"delete", "realms.keycloak-operator.webhippie.de", realmCRName, "-n", testNamespace, "--ignore-not-found", "--wait=false"},
			{"delete", "keycloaks.keycloak-operator.webhippie.de", keycloakCRName, "-n", testNamespace, "--ignore-not-found", "--wait=false"},
		} {
			cmd := exec.Command("kubectl", args...)
			_, _ = utils.Run(cmd)
		}
	})

	// ── Keycloak CR ──────────────────────────────────────────────────────────

	Context("Keycloak CR", func() {
		It("creates a Keycloak CR and becomes connected", func() {
			manifest := fmt.Sprintf(`
apiVersion: keycloak-operator.webhippie.de/v1alpha1
kind: Keycloak
metadata:
  name: %s
  namespace: %s
spec:
  url: %s
  username:
    value: admin
  password:
    value: admin
  insecureSkipVerify: true
`, keycloakCRName, testNamespace, keycloakServiceURL)

			cmd := exec.Command("kubectl", "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(manifest)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to apply Keycloak CR")

			By("waiting for status.connected to become true")
			Eventually(func(g Gomega) {
				out, err := kubectlGetJSONPath(
					"keycloaks.keycloak-operator.webhippie.de", keycloakCRName, testNamespace,
					"{.status.connected}",
				)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(out).To(Equal("true"), "Keycloak CR should be connected")
			}, reconcileTimeout, reconcilePolling).Should(Succeed())
		})
	})

	// ── Realm CR ─────────────────────────────────────────────────────────────

	Context("Realm CR", func() {
		It("creates a Realm CR and populates status.keycloakID", func() {
			manifest := fmt.Sprintf(`
apiVersion: keycloak-operator.webhippie.de/v1alpha1
kind: Realm
metadata:
  name: %s
  namespace: %s
spec:
  keycloakRef:
    name: %s
  realmName: %s
  displayName: E2E Test Realm
`, realmCRName, testNamespace, keycloakCRName, realmName)

			cmd := exec.Command("kubectl", "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(manifest)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to apply Realm CR")

			By("waiting for status.keycloakID to be populated")
			Eventually(func(g Gomega) {
				out, err := kubectlGetJSONPath(
					"realms.keycloak-operator.webhippie.de", realmCRName, testNamespace,
					"{.status.keycloakID}",
				)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(out).NotTo(BeEmpty(), "Realm status.keycloakID should be set")
			}, reconcileTimeout, reconcilePolling).Should(Succeed())
		})

		It("removes the realm from Keycloak when the CR is deleted", func() {
			cmd := exec.Command("kubectl", "delete",
				"realms.keycloak-operator.webhippie.de", realmCRName,
				"-n", testNamespace,
				"--timeout=2m",
			)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to delete Realm CR")

			By("confirming the CR is gone")
			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get",
					"realms.keycloak-operator.webhippie.de", realmCRName,
					"-n", testNamespace,
				)
				out, _ := utils.Run(cmd)
				g.Expect(out).To(ContainSubstring("not found"))
			}, reconcileTimeout, reconcilePolling).Should(Succeed())
		})

		It("re-creates the Realm CR for subsequent tests", func() {
			manifest := fmt.Sprintf(`
apiVersion: keycloak-operator.webhippie.de/v1alpha1
kind: Realm
metadata:
  name: %s
  namespace: %s
spec:
  keycloakRef:
    name: %s
  realmName: %s
  displayName: E2E Test Realm
`, realmCRName, testNamespace, keycloakCRName, realmName)

			cmd := exec.Command("kubectl", "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(manifest)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to re-apply Realm CR")

			Eventually(func(g Gomega) {
				out, err := kubectlGetJSONPath(
					"realms.keycloak-operator.webhippie.de", realmCRName, testNamespace,
					"{.status.keycloakID}",
				)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(out).NotTo(BeEmpty())
			}, reconcileTimeout, reconcilePolling).Should(Succeed())
		})
	})

	// ── User CR ──────────────────────────────────────────────────────────────

	Context("User CR", func() {
		It("creates a User CR and populates status.keycloakID", func() {
			enabled := true
			_ = enabled
			manifest := fmt.Sprintf(`
apiVersion: keycloak-operator.webhippie.de/v1alpha1
kind: User
metadata:
  name: %s
  namespace: %s
spec:
  realmRef:
    name: %s
  username: e2e-testuser
  enabled: true
  firstName: E2E
  lastName: TestUser
  email: e2e-testuser@example.com
`, userCRName, testNamespace, realmCRName)

			cmd := exec.Command("kubectl", "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(manifest)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to apply User CR")

			By("waiting for status.keycloakID to be populated")
			Eventually(func(g Gomega) {
				out, err := kubectlGetJSONPath(
					"users.keycloak-operator.webhippie.de", userCRName, testNamespace,
					"{.status.keycloakID}",
				)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(out).NotTo(BeEmpty(), "User status.keycloakID should be set")
			}, reconcileTimeout, reconcilePolling).Should(Succeed())
		})

		It("removes the user from Keycloak when the CR is deleted", func() {
			cmd := exec.Command("kubectl", "delete",
				"users.keycloak-operator.webhippie.de", userCRName,
				"-n", testNamespace,
				"--timeout=2m",
			)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to delete User CR")

			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get",
					"users.keycloak-operator.webhippie.de", userCRName,
					"-n", testNamespace,
				)
				out, _ := utils.Run(cmd)
				g.Expect(out).To(ContainSubstring("not found"))
			}, reconcileTimeout, reconcilePolling).Should(Succeed())
		})
	})

	// ── Group CR ─────────────────────────────────────────────────────────────

	Context("Group CR", func() {
		It("creates a Group CR and populates status.keycloakID", func() {
			manifest := fmt.Sprintf(`
apiVersion: keycloak-operator.webhippie.de/v1alpha1
kind: Group
metadata:
  name: %s
  namespace: %s
spec:
  realmRef:
    name: %s
  name: e2e-testgroup
`, groupCRName, testNamespace, realmCRName)

			cmd := exec.Command("kubectl", "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(manifest)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to apply Group CR")

			By("waiting for status.keycloakID to be populated")
			Eventually(func(g Gomega) {
				out, err := kubectlGetJSONPath(
					"groups.keycloak-operator.webhippie.de", groupCRName, testNamespace,
					"{.status.keycloakID}",
				)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(out).NotTo(BeEmpty(), "Group status.keycloakID should be set")
			}, reconcileTimeout, reconcilePolling).Should(Succeed())
		})

		It("removes the group from Keycloak when the CR is deleted", func() {
			cmd := exec.Command("kubectl", "delete",
				"groups.keycloak-operator.webhippie.de", groupCRName,
				"-n", testNamespace,
				"--timeout=2m",
			)
			_, err := utils.Run(cmd)
			Expect(err).NotTo(HaveOccurred(), "Failed to delete Group CR")

			Eventually(func(g Gomega) {
				cmd := exec.Command("kubectl", "get",
					"groups.keycloak-operator.webhippie.de", groupCRName,
					"-n", testNamespace,
				)
				out, _ := utils.Run(cmd)
				g.Expect(out).To(ContainSubstring("not found"))
			}, reconcileTimeout, reconcilePolling).Should(Succeed())
		})
	})
})

// kubectlGetJSONPath returns the JSONPath value for a named resource in a namespace.
func kubectlGetJSONPath(resource, name, ns, jsonpath string) (string, error) {
	cmd := exec.Command("kubectl", "get", resource, name,
		"-n", ns,
		"-o", fmt.Sprintf("jsonpath=%s", jsonpath),
	)
	out, err := utils.Run(cmd)
	return strings.TrimSpace(out), err
}
