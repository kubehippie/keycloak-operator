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

var _ = Describe("groupMembershipProtocolMapperToGocloak", func() {
	It("maps all supported fields", func() {
		mapper := &v1alpha1.GroupMembershipProtocolMapper{
			Spec: v1alpha1.GroupMembershipProtocolMapperSpec{
				Name:                    testGroupName,
				ClaimName:               testGroupName,
				FullPath:                boolPtr(false),
				AddToIDToken:            boolPtr(true),
				AddToAccessToken:        boolPtr(false),
				AddToUserInfo:           boolPtr(true),
				AddToTokenIntrospection: boolPtr(false),
			},
		}

		got := groupMembershipProtocolMapperToGocloak(mapper)
		Expect(*got.Name).To(Equal(testGroupName))
		Expect(*got.Protocol).To(Equal("openid-connect"))
		Expect(*got.ProtocolMapper).To(Equal("oidc-group-membership-mapper"))
		Expect(got.Config).To(HaveKeyWithValue("full.path", "false"))
		Expect(got.Config).To(HaveKeyWithValue("claim.name", testGroupName))
		Expect(got.Config).To(HaveKeyWithValue("id.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("access.token.claim", "false"))
		Expect(got.Config).To(HaveKeyWithValue("userinfo.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("introspection.token.claim", "false"))
	})

	It("applies true defaults for nil booleans", func() {
		mapper := &v1alpha1.GroupMembershipProtocolMapper{
			Spec: v1alpha1.GroupMembershipProtocolMapperSpec{Name: testGroupName, ClaimName: testGroupName},
		}

		got := groupMembershipProtocolMapperToGocloak(mapper)
		Expect(got.Config).To(HaveKeyWithValue("full.path", "true"))
		Expect(got.Config).To(HaveKeyWithValue("id.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("access.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("userinfo.token.claim", "true"))
		Expect(got.Config).To(HaveKeyWithValue("introspection.token.claim", "true"))
	})
})
