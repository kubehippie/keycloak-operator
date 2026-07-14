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
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("audienceProtocolMapperToGocloak", func() {
	It("maps included client audience", func() {
		mapper := &v1alpha1.AudienceProtocolMapper{
			Spec: v1alpha1.AudienceProtocolMapperSpec{
				Name:                    "audience",
				IncludedClientAudience:  stringPtr("account"),
				AddToIDToken:            boolPtr(false),
				AddToAccessToken:        boolPtr(true),
				AddToTokenIntrospection: boolPtr(false),
			},
		}

		got := audienceProtocolMapperToGocloak(mapper)
		Expect(*got.ProtocolMapper).To(Equal("oidc-audience-mapper"))
		Expect(got.Config).To(HaveKeyWithValue("included.client.audience", "account"))
		Expect(got.Config).NotTo(HaveKey("included.custom.audience"))
		Expect(got.Config).To(HaveKeyWithValue("id.token.claim", "false"))
		Expect(got.Config).To(HaveKeyWithValue("access.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("introspection.token.claim", "false"))
	})

	It("maps included custom audience and defaults booleans", func() {
		mapper := &v1alpha1.AudienceProtocolMapper{
			Spec: v1alpha1.AudienceProtocolMapperSpec{
				Name:                   "audience",
				IncludedCustomAudience: stringPtr("my-api"),
			},
		}

		got := audienceProtocolMapperToGocloak(mapper)
		Expect(got.Config).To(HaveKeyWithValue("included.custom.audience", "my-api"))
		Expect(got.Config).NotTo(HaveKey("included.client.audience"))
		Expect(got.Config).To(HaveKeyWithValue("id.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("access.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("introspection.token.claim", "true"))
	})
})
