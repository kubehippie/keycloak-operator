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

var _ = Describe("userAttributeProtocolMapperToGocloak", func() {
	It("maps all supported fields", func() {
		mapper := &v1alpha1.UserAttributeProtocolMapper{
			Spec: v1alpha1.UserAttributeProtocolMapperSpec{
				Name:                    testAttributeDepartment,
				UserAttribute:           testAttributeDepartment,
				ClaimName:               testAttributeDepartment,
				ClaimValueType:          stringPtr("long"),
				Multivalued:             boolPtr(false),
				AggregateAttributes:     boolPtr(true),
				AddToIDToken:            boolPtr(false),
				AddToAccessToken:        boolPtr(true),
				AddToUserInfo:           boolPtr(false),
				AddToTokenIntrospection: boolPtr(true),
			},
		}

		got := userAttributeProtocolMapperToGocloak(mapper)
		Expect(*got.ProtocolMapper).To(Equal("oidc-usermodel-attribute-mapper"))
		Expect(got.Config).To(HaveKeyWithValue("user.attribute", testAttributeDepartment))
		Expect(got.Config).To(HaveKeyWithValue("claim.name", testAttributeDepartment))
		Expect(got.Config).To(HaveKeyWithValue("jsonType.label", "long"))
		Expect(got.Config).To(HaveKeyWithValue("multivalued", "false"))
		Expect(got.Config).To(HaveKeyWithValue("aggregate.attrs", "true"))
		Expect(got.Config).To(HaveKeyWithValue("id.token.claim", "false"))
		Expect(got.Config).To(HaveKeyWithValue("access.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("userinfo.token.claim", "false"))
		Expect(got.Config).To(HaveKeyWithValue("introspection.token.claim", "true"))
	})

	It("defaults claimValueType and all booleans", func() {
		mapper := &v1alpha1.UserAttributeProtocolMapper{
			Spec: v1alpha1.UserAttributeProtocolMapperSpec{Name: testAttributeTeam, UserAttribute: testAttributeTeam, ClaimName: testAttributeTeam},
		}

		got := userAttributeProtocolMapperToGocloak(mapper)
		Expect(got.Config).To(HaveKeyWithValue("jsonType.label", "String"))
		Expect(got.Config).To(HaveKeyWithValue("multivalued", "true"))
		Expect(got.Config).To(HaveKeyWithValue("aggregate.attrs", "true"))
		Expect(got.Config).To(HaveKeyWithValue("id.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("access.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("userinfo.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("introspection.token.claim", "true"))
	})
})
