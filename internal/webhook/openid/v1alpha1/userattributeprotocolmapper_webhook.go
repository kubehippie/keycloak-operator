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

	openidv1alpha1 "github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var userattributeprotocolmapperlog = logf.Log.WithName("userattributeprotocolmapper-resource")

func SetupUserAttributeProtocolMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &openidv1alpha1.UserAttributeProtocolMapper{}).
		WithValidator(&UserAttributeProtocolMapperCustomValidator{}).
		WithDefaulter(&UserAttributeProtocolMapperCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-openid-keycloak-operator-webhippie-de-v1alpha1-userattributeprotocolmapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=userattributeprotocolmappers,verbs=create;update,versions=v1alpha1,name=muserattributeprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

type UserAttributeProtocolMapperCustomDefaulter struct{}

var _ admission.Defaulter[*openidv1alpha1.UserAttributeProtocolMapper] = &UserAttributeProtocolMapperCustomDefaulter{}

func (d *UserAttributeProtocolMapperCustomDefaulter) Default(_ context.Context, mapper *openidv1alpha1.UserAttributeProtocolMapper) error {
	userattributeprotocolmapperlog.Info("Defaulting for UserAttributeProtocolMapper", "name", mapper.GetName())

	if mapper.Spec.ClaimValueType == nil {
		mapper.Spec.ClaimValueType = stringPtr("String")
	}
	setDefaultTrue(&mapper.Spec.Multivalued)
	setDefaultTrue(&mapper.Spec.AggregateAttributes)
	setDefaultTrue(&mapper.Spec.AddToIDToken)
	setDefaultTrue(&mapper.Spec.AddToAccessToken)
	setDefaultTrue(&mapper.Spec.AddToUserInfo)
	setDefaultTrue(&mapper.Spec.AddToTokenIntrospection)

	return nil
}

// +kubebuilder:webhook:path=/validate-openid-keycloak-operator-webhippie-de-v1alpha1-userattributeprotocolmapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=userattributeprotocolmappers,verbs=create;update,versions=v1alpha1,name=vuserattributeprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

type UserAttributeProtocolMapperCustomValidator struct{}

var _ admission.Validator[*openidv1alpha1.UserAttributeProtocolMapper] = &UserAttributeProtocolMapperCustomValidator{}

func (v *UserAttributeProtocolMapperCustomValidator) ValidateCreate(_ context.Context, mapper *openidv1alpha1.UserAttributeProtocolMapper) (admission.Warnings, error) {
	userattributeprotocolmapperlog.Info("Validation for UserAttributeProtocolMapper upon creation", "name", mapper.GetName())

	if errs := validateUserAttributeProtocolMapper(mapper); len(errs) > 0 {
		return nil, errs.ToAggregate()
	}

	return nil, nil
}

func (v *UserAttributeProtocolMapperCustomValidator) ValidateUpdate(_ context.Context, oldMapper, mapper *openidv1alpha1.UserAttributeProtocolMapper) (admission.Warnings, error) {
	userattributeprotocolmapperlog.Info("Validation for UserAttributeProtocolMapper upon update", "name", mapper.GetName())

	var allErrs field.ErrorList
	allErrs = append(allErrs, validateUserAttributeProtocolMapper(mapper)...)

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

func (v *UserAttributeProtocolMapperCustomValidator) ValidateDelete(_ context.Context, mapper *openidv1alpha1.UserAttributeProtocolMapper) (admission.Warnings, error) {
	userattributeprotocolmapperlog.Info("Validation for UserAttributeProtocolMapper upon deletion", "name", mapper.GetName())
	return nil, nil
}

func validateUserAttributeProtocolMapper(mapper *openidv1alpha1.UserAttributeProtocolMapper) field.ErrorList {
	errs := make(field.ErrorList, 0, 5)

	errs = append(errs, validateClientRef(mapper.Spec.ClientRef, field.NewPath("spec", "clientRef"))...)
	errs = append(errs, validateRequiredString(mapper.Spec.Name, field.NewPath("spec", "name"), "name is required")...)
	errs = append(errs, validateRequiredString(mapper.Spec.UserAttribute, field.NewPath("spec", "userAttribute"), "userAttribute is required")...)
	errs = append(errs, validateRequiredString(mapper.Spec.ClaimName, field.NewPath("spec", "claimName"), "claimName is required")...)
	errs = append(errs, validateClaimValueType(mapper.Spec.ClaimValueType, field.NewPath("spec", "claimValueType"))...)

	return errs
}
