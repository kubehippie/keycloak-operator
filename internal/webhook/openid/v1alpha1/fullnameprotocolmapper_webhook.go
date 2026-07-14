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

var fullnameprotocolmapperlog = logf.Log.WithName("fullnameprotocolmapper-resource")

func SetupFullNameProtocolMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &openidv1alpha1.FullNameProtocolMapper{}).
		WithValidator(&FullNameProtocolMapperCustomValidator{}).
		WithDefaulter(&FullNameProtocolMapperCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-openid-keycloak-operator-webhippie-de-v1alpha1-fullnameprotocolmapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=fullnameprotocolmappers,verbs=create;update,versions=v1alpha1,name=mfullnameprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

type FullNameProtocolMapperCustomDefaulter struct{}

var _ admission.Defaulter[*openidv1alpha1.FullNameProtocolMapper] = &FullNameProtocolMapperCustomDefaulter{}

func (d *FullNameProtocolMapperCustomDefaulter) Default(_ context.Context, mapper *openidv1alpha1.FullNameProtocolMapper) error {
	fullnameprotocolmapperlog.Info("Defaulting for FullNameProtocolMapper", "name", mapper.GetName())

	setDefaultTrue(&mapper.Spec.AddToIDToken)
	setDefaultTrue(&mapper.Spec.AddToAccessToken)
	setDefaultTrue(&mapper.Spec.AddToUserInfo)

	return nil
}

// +kubebuilder:webhook:path=/validate-openid-keycloak-operator-webhippie-de-v1alpha1-fullnameprotocolmapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=fullnameprotocolmappers,verbs=create;update,versions=v1alpha1,name=vfullnameprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

type FullNameProtocolMapperCustomValidator struct{}

var _ admission.Validator[*openidv1alpha1.FullNameProtocolMapper] = &FullNameProtocolMapperCustomValidator{}

func (v *FullNameProtocolMapperCustomValidator) ValidateCreate(_ context.Context, mapper *openidv1alpha1.FullNameProtocolMapper) (admission.Warnings, error) {
	fullnameprotocolmapperlog.Info("Validation for FullNameProtocolMapper upon creation", "name", mapper.GetName())

	if errs := validateFullNameProtocolMapper(mapper); len(errs) > 0 {
		return nil, errs.ToAggregate()
	}

	return nil, nil
}

func (v *FullNameProtocolMapperCustomValidator) ValidateUpdate(_ context.Context, oldMapper, mapper *openidv1alpha1.FullNameProtocolMapper) (admission.Warnings, error) {
	fullnameprotocolmapperlog.Info("Validation for FullNameProtocolMapper upon update", "name", mapper.GetName())

	var allErrs field.ErrorList
	allErrs = append(allErrs, validateFullNameProtocolMapper(mapper)...)

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

func (v *FullNameProtocolMapperCustomValidator) ValidateDelete(_ context.Context, mapper *openidv1alpha1.FullNameProtocolMapper) (admission.Warnings, error) {
	fullnameprotocolmapperlog.Info("Validation for FullNameProtocolMapper upon deletion", "name", mapper.GetName())
	return nil, nil
}

func validateFullNameProtocolMapper(mapper *openidv1alpha1.FullNameProtocolMapper) field.ErrorList {
	errs := make(field.ErrorList, 0, 2)

	errs = append(errs, validateClientRef(mapper.Spec.ClientRef, field.NewPath("spec", "clientRef"))...)
	errs = append(errs, validateRequiredString(mapper.Spec.Name, field.NewPath("spec", "name"), "name is required")...)

	return errs
}
