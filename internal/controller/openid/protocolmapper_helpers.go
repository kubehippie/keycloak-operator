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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Nerzal/gocloak/v14"
)

const (
	openidConnectProtocol = "openid-connect"
	claimValueTypeString  = "String"

	configKeyIDTokenClaim            = "id.token.claim"
	configKeyAccessTokenClaim        = "access.token.claim"
	configKeyUserinfoTokenClaim      = "userinfo.token.claim"
	configKeyIntrospectionTokenClaim = "introspection.token.claim"
	configKeyClaimName               = "claim.name"
)

func boolStringDefaultTrue(value *bool) string {
	if value == nil {
		return "true"
	}

	return strconv.FormatBool(*value)
}

func stringValueDefault(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fallback
	}

	return *value
}

func setOptionalConfig(config map[string]string, key string, value *string) {
	if value == nil {
		return
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return
	}

	config[key] = trimmed
}

func findClientProtocolMapperByName(client *gocloak.Client, name string) *gocloak.ProtocolMapperRepresentation {
	if client == nil {
		return nil
	}

	for i := range client.ProtocolMappers {
		mapper := &client.ProtocolMappers[i]
		if mapper.Name != nil && *mapper.Name == name && mapper.ID != nil {
			return mapper
		}
	}

	return nil
}

func isNotFoundAPIError(err error) bool {
	var apiErr *gocloak.APIError
	return errors.As(err, &apiErr) && apiErr.Code == 404
}

func resolveDefaultScopesPlan(realmScopeIDs map[string]string, desiredNames []string, currentScopes []*gocloak.ClientScope) ([]string, []string, error) {
	desiredSet := make(map[string]struct{}, len(desiredNames))
	toAdd := make([]string, 0)
	for _, desiredName := range desiredNames {
		trimmed := strings.TrimSpace(desiredName)
		if trimmed == "" {
			continue
		}
		if _, ok := desiredSet[trimmed]; ok {
			continue
		}

		scopeID, ok := realmScopeIDs[trimmed]
		if !ok {
			return nil, nil, fmt.Errorf("default client scope %q does not exist in realm", trimmed)
		}

		desiredSet[trimmed] = struct{}{}
		toAdd = append(toAdd, scopeID)
	}

	currentByName := make(map[string]string, len(currentScopes))
	for _, current := range currentScopes {
		if current == nil || current.Name == nil {
			continue
		}

		if current.ID != nil {
			currentByName[*current.Name] = *current.ID
			continue
		}

		if scopeID, ok := realmScopeIDs[*current.Name]; ok {
			currentByName[*current.Name] = scopeID
		}
	}

	filteredToAdd := make([]string, 0, len(toAdd))
	for _, scopeID := range toAdd {
		alreadyAttached := false
		for _, currentID := range currentByName {
			if currentID == scopeID {
				alreadyAttached = true
				break
			}
		}
		if !alreadyAttached {
			filteredToAdd = append(filteredToAdd, scopeID)
		}
	}

	toRemove := make([]string, 0)
	for currentName, currentID := range currentByName {
		if _, ok := desiredSet[currentName]; !ok {
			toRemove = append(toRemove, currentID)
		}
	}

	return filteredToAdd, toRemove, nil
}
