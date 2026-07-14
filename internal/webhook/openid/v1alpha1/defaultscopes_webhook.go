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
	"strings"

	openidv1alpha1 "github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:unused
// log is for logging in this package.
var defaultscopeslog = logf.Log.WithName("defaultscopes-resource")

// SetupDefaultScopesWebhookWithManager registers the webhook for DefaultScopes in the manager.
func SetupDefaultScopesWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &openidv1alpha1.DefaultScopes{}).
		WithValidator(&DefaultScopesCustomValidator{}).
		WithDefaulter(&DefaultScopesCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-openid-keycloak-operator-webhippie-de-v1alpha1-defaultscopes,mutating=true,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=defaultscopes,verbs=create;update,versions=v1alpha1,name=mdefaultscopes-v1alpha1.kb.io,admissionReviewVersions=v1

type DefaultScopesCustomDefaulter struct{}

var _ admission.Defaulter[*openidv1alpha1.DefaultScopes] = &DefaultScopesCustomDefaulter{}

func (d *DefaultScopesCustomDefaulter) Default(_ context.Context, defaultScopes *openidv1alpha1.DefaultScopes) error {
	defaultscopeslog.Info("Defaulting for DefaultScopes", "name", defaultScopes.GetName())
	return nil
}

// +kubebuilder:webhook:path=/validate-openid-keycloak-operator-webhippie-de-v1alpha1-defaultscopes,mutating=false,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=defaultscopes,verbs=create;update,versions=v1alpha1,name=vdefaultscopes-v1alpha1.kb.io,admissionReviewVersions=v1

type DefaultScopesCustomValidator struct{}

var _ admission.Validator[*openidv1alpha1.DefaultScopes] = &DefaultScopesCustomValidator{}

func (v *DefaultScopesCustomValidator) ValidateCreate(_ context.Context, defaultScopes *openidv1alpha1.DefaultScopes) (admission.Warnings, error) {
	defaultscopeslog.Info("Validation for DefaultScopes upon creation", "name", defaultScopes.GetName())

	if errs := validateDefaultScopes(defaultScopes); len(errs) > 0 {
		return nil, errs.ToAggregate()
	}

	return nil, nil
}

func (v *DefaultScopesCustomValidator) ValidateUpdate(_ context.Context, oldDefaultScopes, defaultScopes *openidv1alpha1.DefaultScopes) (admission.Warnings, error) {
	defaultscopeslog.Info("Validation for DefaultScopes upon update", "name", defaultScopes.GetName())

	var allErrs field.ErrorList
	allErrs = append(allErrs, validateDefaultScopes(defaultScopes)...)

	if !clientRefEqual(oldDefaultScopes.Spec.ClientRef, defaultScopes.Spec.ClientRef) {
		allErrs = append(allErrs, field.Forbidden(
			field.NewPath("spec", "clientRef"),
			"clientRef is immutable and cannot be changed after creation",
		))
	}

	if len(allErrs) > 0 {
		return nil, allErrs.ToAggregate()
	}

	return nil, nil
}

func (v *DefaultScopesCustomValidator) ValidateDelete(_ context.Context, defaultScopes *openidv1alpha1.DefaultScopes) (admission.Warnings, error) {
	defaultscopeslog.Info("Validation for DefaultScopes upon deletion", "name", defaultScopes.GetName())
	return nil, nil
}

func validateDefaultScopes(defaultScopes *openidv1alpha1.DefaultScopes) field.ErrorList {
	var errs field.ErrorList

	errs = append(errs, validateClientRef(defaultScopes.Spec.ClientRef, field.NewPath("spec", "clientRef"))...)

	if len(defaultScopes.Spec.DefaultScopes) == 0 {
		errs = append(errs, field.Required(
			field.NewPath("spec", "defaultScopes"),
			"defaultScopes must contain at least one client scope name",
		))
	}

	for i, scopeName := range defaultScopes.Spec.DefaultScopes {
		if strings.TrimSpace(scopeName) == "" {
			errs = append(errs, field.Required(
				field.NewPath("spec", "defaultScopes").Index(i),
				"default scope name must not be empty",
			))
		}
	}

	return errs
}
