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
	"github.com/Nerzal/gocloak/v14"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("resolveDefaultScopesPlan", func() {
	It("returns scopes to add and remove by diffing names", func() {
		realmScopes := map[string]string{
			testScopeEmail: "scope-email",
			"profile":      "scope-profile",
			"roles":        "scope-roles",
		}
		currentScopes := []*gocloak.ClientScope{
			{Name: gocloak.StringP(testScopeEmail), ID: gocloak.StringP("scope-email")},
			{Name: gocloak.StringP("roles"), ID: gocloak.StringP("scope-roles")},
		}

		toAdd, toRemove, err := resolveDefaultScopesPlan(realmScopes, []string{testScopeEmail, "profile"}, currentScopes)
		Expect(err).NotTo(HaveOccurred())
		Expect(toAdd).To(ConsistOf("scope-profile"))
		Expect(toRemove).To(ConsistOf("scope-roles"))
	})

	It("ignores duplicate desired scope names", func() {
		realmScopes := map[string]string{testScopeEmail: "scope-email"}

		toAdd, toRemove, err := resolveDefaultScopesPlan(realmScopes, []string{testScopeEmail, testScopeEmail}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(toAdd).To(ConsistOf("scope-email"))
		Expect(toRemove).To(BeEmpty())
	})

	It("fails when a desired scope does not exist in the realm", func() {
		_, _, err := resolveDefaultScopesPlan(map[string]string{}, []string{"missing"}, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("missing"))
	})
})
