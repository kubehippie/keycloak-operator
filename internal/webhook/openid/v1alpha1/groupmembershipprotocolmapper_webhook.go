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

var groupmembershipprotocolmapperlog = logf.Log.WithName("groupmembershipprotocolmapper-resource")

func SetupGroupMembershipProtocolMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &openidv1alpha1.GroupMembershipProtocolMapper{}).
		WithValidator(&GroupMembershipProtocolMapperCustomValidator{}).
		WithDefaulter(&GroupMembershipProtocolMapperCustomDefaulter{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-openid-keycloak-operator-webhippie-de-v1alpha1-groupmembershipprotocolmapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=groupmembershipprotocolmappers,verbs=create;update,versions=v1alpha1,name=mgroupmembershipprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

type GroupMembershipProtocolMapperCustomDefaulter struct{}

var _ admission.Defaulter[*openidv1alpha1.GroupMembershipProtocolMapper] = &GroupMembershipProtocolMapperCustomDefaulter{}

func (d *GroupMembershipProtocolMapperCustomDefaulter) Default(_ context.Context, mapper *openidv1alpha1.GroupMembershipProtocolMapper) error {
	groupmembershipprotocolmapperlog.Info("Defaulting for GroupMembershipProtocolMapper", "name", mapper.GetName())

	setDefaultTrue(&mapper.Spec.FullPath)
	setDefaultTrue(&mapper.Spec.AddToIDToken)
	setDefaultTrue(&mapper.Spec.AddToAccessToken)
	setDefaultTrue(&mapper.Spec.AddToUserInfo)
	setDefaultTrue(&mapper.Spec.AddToTokenIntrospection)

	return nil
}

// +kubebuilder:webhook:path=/validate-openid-keycloak-operator-webhippie-de-v1alpha1-groupmembershipprotocolmapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=groupmembershipprotocolmappers,verbs=create;update,versions=v1alpha1,name=vgroupmembershipprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

type GroupMembershipProtocolMapperCustomValidator struct{}

var _ admission.Validator[*openidv1alpha1.GroupMembershipProtocolMapper] = &GroupMembershipProtocolMapperCustomValidator{}

func (v *GroupMembershipProtocolMapperCustomValidator) ValidateCreate(_ context.Context, mapper *openidv1alpha1.GroupMembershipProtocolMapper) (admission.Warnings, error) {
	groupmembershipprotocolmapperlog.Info("Validation for GroupMembershipProtocolMapper upon creation", "name", mapper.GetName())

	if errs := validateGroupMembershipProtocolMapper(mapper); len(errs) > 0 {
		return nil, errs.ToAggregate()
	}

	return nil, nil
}

func (v *GroupMembershipProtocolMapperCustomValidator) ValidateUpdate(_ context.Context, oldMapper, mapper *openidv1alpha1.GroupMembershipProtocolMapper) (admission.Warnings, error) {
	groupmembershipprotocolmapperlog.Info("Validation for GroupMembershipProtocolMapper upon update", "name", mapper.GetName())

	var allErrs field.ErrorList
	allErrs = append(allErrs, validateGroupMembershipProtocolMapper(mapper)...)

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

func (v *GroupMembershipProtocolMapperCustomValidator) ValidateDelete(_ context.Context, mapper *openidv1alpha1.GroupMembershipProtocolMapper) (admission.Warnings, error) {
	groupmembershipprotocolmapperlog.Info("Validation for GroupMembershipProtocolMapper upon deletion", "name", mapper.GetName())
	return nil, nil
}

func validateGroupMembershipProtocolMapper(mapper *openidv1alpha1.GroupMembershipProtocolMapper) field.ErrorList {
	errs := make(field.ErrorList, 0, 3)

	errs = append(errs, validateClientRef(mapper.Spec.ClientRef, field.NewPath("spec", "clientRef"))...)
	errs = append(errs, validateRequiredString(mapper.Spec.Name, field.NewPath("spec", "name"), "name is required")...)
	errs = append(errs, validateRequiredString(mapper.Spec.ClaimName, field.NewPath("spec", "claimName"), "claimName is required")...)

	return errs
}
