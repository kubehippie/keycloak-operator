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
	"context"
	"reflect"
	"strings"

	identityv1alpha1 "github.com/kubehippie/keycloak-operator/api/identity/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:unused
// log is for logging in this package.
var oidcadvancedgroupmapperlog = logf.Log.WithName("oidcadvancedgroupmapper-resource")

// SetupOIDCAdvancedGroupMapperWebhookWithManager registers the webhook for OIDCAdvancedGroupMapper in the manager.
func SetupOIDCAdvancedGroupMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &identityv1alpha1.OIDCAdvancedGroupMapper{}).
		WithValidator(&OIDCAdvancedGroupMapperCustomValidator{}).
		WithDefaulter(&OIDCAdvancedGroupMapperCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-identity-keycloak-operator-webhippie-de-v1alpha1-oidcadvancedgroupmapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=oidcadvancedgroupmappers,verbs=create;update,versions=v1alpha1,name=moidcadvancedgroupmapper-v1alpha1.kb.io,admissionReviewVersions=v1

// OIDCAdvancedGroupMapperCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind OIDCAdvancedGroupMapper when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type OIDCAdvancedGroupMapperCustomDefaulter struct{}

var _ admission.Defaulter[*identityv1alpha1.OIDCAdvancedGroupMapper] = &OIDCAdvancedGroupMapperCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind OIDCAdvancedGroupMapper.
func (d *OIDCAdvancedGroupMapperCustomDefaulter) Default(_ context.Context, mapper *identityv1alpha1.OIDCAdvancedGroupMapper) error {
	oidcadvancedgroupmapperlog.Info("Defaulting for OIDCAdvancedGroupMapper", "name", mapper.GetName())

	if mapper.Spec.ClaimValuesRegex == nil {
		mapper.Spec.ClaimValuesRegex = new(bool)
	}

	return nil
}

// +kubebuilder:webhook:path=/validate-identity-keycloak-operator-webhippie-de-v1alpha1-oidcadvancedgroupmapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=oidcadvancedgroupmappers,verbs=create;update,versions=v1alpha1,name=voidcadvancedgroupmapper-v1alpha1.kb.io,admissionReviewVersions=v1

// OIDCAdvancedGroupMapperCustomValidator struct is responsible for validating the OIDCAdvancedGroupMapper resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type OIDCAdvancedGroupMapperCustomValidator struct{}

var _ admission.Validator[*identityv1alpha1.OIDCAdvancedGroupMapper] = &OIDCAdvancedGroupMapperCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type OIDCAdvancedGroupMapper.
func (v *OIDCAdvancedGroupMapperCustomValidator) ValidateCreate(_ context.Context, mapper *identityv1alpha1.OIDCAdvancedGroupMapper) (admission.Warnings, error) {
	oidcadvancedgroupmapperlog.Info("Validation for OIDCAdvancedGroupMapper upon creation", "name", mapper.GetName())

	if errs := validateOIDCAdvancedGroupMapper(mapper); len(errs) > 0 {
		return nil, errs.ToAggregate()
	}

	return nil, nil
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type OIDCAdvancedGroupMapper.
func (v *OIDCAdvancedGroupMapperCustomValidator) ValidateUpdate(_ context.Context, oldMapper, mapper *identityv1alpha1.OIDCAdvancedGroupMapper) (admission.Warnings, error) {
	oidcadvancedgroupmapperlog.Info("Validation for OIDCAdvancedGroupMapper upon update", "name", mapper.GetName())

	allErrs := validateOIDCAdvancedGroupMapper(mapper)

	if !reflect.DeepEqual(oldMapper.Spec.IdentityProviderRef, mapper.Spec.IdentityProviderRef) {
		allErrs = append(allErrs, field.Forbidden(
			field.NewPath("spec", "identityProviderRef"),
			"identityProviderRef is immutable and cannot be changed after creation",
		))
	}

	if len(allErrs) > 0 {
		return nil, allErrs.ToAggregate()
	}

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type OIDCAdvancedGroupMapper.
func (v *OIDCAdvancedGroupMapperCustomValidator) ValidateDelete(_ context.Context, mapper *identityv1alpha1.OIDCAdvancedGroupMapper) (admission.Warnings, error) {
	oidcadvancedgroupmapperlog.Info("Validation for OIDCAdvancedGroupMapper upon deletion", "name", mapper.GetName())
	return nil, nil
}

func validateOIDCAdvancedGroupMapper(mapper *identityv1alpha1.OIDCAdvancedGroupMapper) field.ErrorList {
	errs := make(field.ErrorList, 0, 4)

	if mapper.Spec.IdentityProviderRef == nil {
		errs = append(errs, field.Required(
			field.NewPath("spec", "identityProviderRef"),
			"identityProviderRef is required",
		))
	} else if strings.TrimSpace(mapper.Spec.IdentityProviderRef.Name) == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "identityProviderRef", "name"),
			"identityProviderRef.name is required",
		))
	}

	if strings.TrimSpace(mapper.Spec.Name) == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "name"),
			"name is required",
		))
	}

	if len(mapper.Spec.Claims) == 0 {
		errs = append(errs, field.Required(
			field.NewPath("spec", "claims"),
			"claims must contain at least one entry",
		))
	}

	if strings.TrimSpace(mapper.Spec.Group) == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "group"),
			"group is required",
		))
	}

	return errs
}
