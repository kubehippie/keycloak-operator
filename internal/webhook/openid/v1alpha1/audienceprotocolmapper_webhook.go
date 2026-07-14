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

var audienceprotocolmapperlog = logf.Log.WithName("audienceprotocolmapper-resource")

func SetupAudienceProtocolMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &openidv1alpha1.AudienceProtocolMapper{}).
		WithValidator(&AudienceProtocolMapperCustomValidator{}).
		WithDefaulter(&AudienceProtocolMapperCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-openid-keycloak-operator-webhippie-de-v1alpha1-audienceprotocolmapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=audienceprotocolmappers,verbs=create;update,versions=v1alpha1,name=maudienceprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

type AudienceProtocolMapperCustomDefaulter struct{}

var _ admission.Defaulter[*openidv1alpha1.AudienceProtocolMapper] = &AudienceProtocolMapperCustomDefaulter{}

func (d *AudienceProtocolMapperCustomDefaulter) Default(_ context.Context, mapper *openidv1alpha1.AudienceProtocolMapper) error {
	audienceprotocolmapperlog.Info("Defaulting for AudienceProtocolMapper", "name", mapper.GetName())

	setDefaultTrue(&mapper.Spec.AddToIDToken)
	setDefaultTrue(&mapper.Spec.AddToAccessToken)
	setDefaultTrue(&mapper.Spec.AddToTokenIntrospection)

	return nil
}

// +kubebuilder:webhook:path=/validate-openid-keycloak-operator-webhippie-de-v1alpha1-audienceprotocolmapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=audienceprotocolmappers,verbs=create;update,versions=v1alpha1,name=vaudienceprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

type AudienceProtocolMapperCustomValidator struct{}

var _ admission.Validator[*openidv1alpha1.AudienceProtocolMapper] = &AudienceProtocolMapperCustomValidator{}

func (v *AudienceProtocolMapperCustomValidator) ValidateCreate(_ context.Context, mapper *openidv1alpha1.AudienceProtocolMapper) (admission.Warnings, error) {
	audienceprotocolmapperlog.Info("Validation for AudienceProtocolMapper upon creation", "name", mapper.GetName())

	if errs := validateAudienceProtocolMapper(mapper); len(errs) > 0 {
		return nil, errs.ToAggregate()
	}

	return nil, nil
}

func (v *AudienceProtocolMapperCustomValidator) ValidateUpdate(_ context.Context, oldMapper, mapper *openidv1alpha1.AudienceProtocolMapper) (admission.Warnings, error) {
	audienceprotocolmapperlog.Info("Validation for AudienceProtocolMapper upon update", "name", mapper.GetName())

	var allErrs field.ErrorList
	allErrs = append(allErrs, validateAudienceProtocolMapper(mapper)...)

	if !clientRefEqual(oldMapper.Spec.ClientRef, mapper.Spec.ClientRef) {
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

func (v *AudienceProtocolMapperCustomValidator) ValidateDelete(_ context.Context, mapper *openidv1alpha1.AudienceProtocolMapper) (admission.Warnings, error) {
	audienceprotocolmapperlog.Info("Validation for AudienceProtocolMapper upon deletion", "name", mapper.GetName())
	return nil, nil
}

func validateAudienceProtocolMapper(mapper *openidv1alpha1.AudienceProtocolMapper) field.ErrorList {
	var errs field.ErrorList

	errs = append(errs, validateClientRef(mapper.Spec.ClientRef, field.NewPath("spec", "clientRef"))...)
	errs = append(errs, validateRequiredString(mapper.Spec.Name, field.NewPath("spec", "name"), "name is required")...)

	hasClientAudience := mapper.Spec.IncludedClientAudience != nil && strings.TrimSpace(*mapper.Spec.IncludedClientAudience) != ""
	hasCustomAudience := mapper.Spec.IncludedCustomAudience != nil && strings.TrimSpace(*mapper.Spec.IncludedCustomAudience) != ""

	if hasClientAudience == hasCustomAudience {
		errs = append(errs, field.Invalid(
			field.NewPath("spec"),
			mapper.Spec,
			"exactly one of includedClientAudience or includedCustomAudience must be set",
		))
	}

	return errs
}
