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

	"github.com/kubehippie/keycloak-operator/api/identity/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:unused
// log is for logging in this package.
var customprovidermapperlog = logf.Log.WithName("customprovidermapper-resource")

// SetupCustomProviderMapperWebhookWithManager registers the webhook for CustomProviderMapper in the manager.
func SetupCustomProviderMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.CustomProviderMapper{}).
		WithValidator(&CustomProviderMapperCustomValidator{}).
		WithDefaulter(&CustomProviderMapperCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-identity-keycloak-operator-webhippie-de-v1alpha1-customprovidermapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=customprovidermappers,verbs=create;update,versions=v1alpha1,name=mcustomprovidermapper-v1alpha1.kb.io,admissionReviewVersions=v1

// CustomProviderMapperCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind CustomProviderMapper when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type CustomProviderMapperCustomDefaulter struct{}

var _ admission.Defaulter[*v1alpha1.CustomProviderMapper] = &CustomProviderMapperCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind CustomProviderMapper.
func (d *CustomProviderMapperCustomDefaulter) Default(_ context.Context, mapper *v1alpha1.CustomProviderMapper) error {
	customprovidermapperlog.Info("Defaulting for CustomProviderMapper", "name", mapper.GetName())
	return nil
}

// +kubebuilder:webhook:path=/validate-identity-keycloak-operator-webhippie-de-v1alpha1-customprovidermapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=identity.keycloak-operator.webhippie.de,resources=customprovidermappers,verbs=create;update,versions=v1alpha1,name=vcustomprovidermapper-v1alpha1.kb.io,admissionReviewVersions=v1

// CustomProviderMapperCustomValidator struct is responsible for validating the CustomProviderMapper resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type CustomProviderMapperCustomValidator struct{}

var _ admission.Validator[*v1alpha1.CustomProviderMapper] = &CustomProviderMapperCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type CustomProviderMapper.
func (v *CustomProviderMapperCustomValidator) ValidateCreate(_ context.Context, mapper *v1alpha1.CustomProviderMapper) (admission.Warnings, error) {
	customprovidermapperlog.Info("Validation for CustomProviderMapper upon creation", "name", mapper.GetName())

	if errs := validateCustomProviderMapper(mapper); len(errs) > 0 {
		return nil, errs.ToAggregate()
	}

	return nil, nil
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type CustomProviderMapper.
func (v *CustomProviderMapperCustomValidator) ValidateUpdate(_ context.Context, oldMapper, mapper *v1alpha1.CustomProviderMapper) (admission.Warnings, error) {
	customprovidermapperlog.Info("Validation for CustomProviderMapper upon update", "name", mapper.GetName())

	var allErrs field.ErrorList

	if errs := validateCustomProviderMapper(mapper); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

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

// ValidateDelete implements admission.Validator so a webhook will be registered for the type CustomProviderMapper.
func (v *CustomProviderMapperCustomValidator) ValidateDelete(_ context.Context, mapper *v1alpha1.CustomProviderMapper) (admission.Warnings, error) {
	customprovidermapperlog.Info("Validation for CustomProviderMapper upon deletion", "name", mapper.GetName())
	return nil, nil
}

func validateCustomProviderMapper(mapper *v1alpha1.CustomProviderMapper) field.ErrorList {
	var errs field.ErrorList

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

	if strings.TrimSpace(mapper.Spec.MapperType) == "" {
		errs = append(errs, field.Required(
			field.NewPath("spec", "mapperType"),
			"mapperType is required",
		))
	}

	return errs
}
