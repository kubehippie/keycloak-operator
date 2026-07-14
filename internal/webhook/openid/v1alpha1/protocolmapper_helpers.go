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

package v1alpha1

import (
	"strings"

	"github.com/kubehippie/keycloak-operator/api/common"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var supportedClaimValueTypes = []string{"String", "JSON", "long", "int", "boolean"}

func boolPtr(value bool) *bool {
	return &value
}

func stringPtr(value string) *string {
	return &value
}

func setDefaultTrue(target **bool) {
	if *target == nil {
		*target = boolPtr(true)
	}
}

func normalizeClientRef(ref *common.ClientRef) (kind, namespace, name string) {
	if ref == nil {
		return "", "", ""
	}

	kind = ref.Kind
	if kind == "" {
		kind = "OpenIDClient"
	}

	return kind, ref.Namespace, ref.Name
}

func clientRefEqual(a, b *common.ClientRef) bool {
	ak, ans, an := normalizeClientRef(a)
	bk, bns, bn := normalizeClientRef(b)
	return ak == bk && ans == bns && an == bn
}

func validateClientRef(ref *common.ClientRef, path *field.Path) field.ErrorList {
	var errs field.ErrorList

	if ref == nil {
		errs = append(errs, field.Required(path, "clientRef is required"))
		return errs
	}

	if strings.TrimSpace(ref.Name) == "" {
		errs = append(errs, field.Required(path.Child("name"), "clientRef.name is required"))
	}

	return errs
}

func validateRequiredString(value string, path *field.Path, message string) field.ErrorList {
	if strings.TrimSpace(value) == "" {
		return field.ErrorList{field.Required(path, message)}
	}

	return nil
}

func validateClaimValueType(value *string, path *field.Path) field.ErrorList {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	for _, allowed := range supportedClaimValueTypes {
		if trimmed == allowed {
			return nil
		}
	}

	return field.ErrorList{field.NotSupported(path, trimmed, supportedClaimValueTypes)}
}
