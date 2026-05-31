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

	"github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// nolint:unused
// log is for logging in this package.
var groupmembershipprotocolmapperlog = logf.Log.WithName("groupmembershipprotocolmapper-resource")

// SetupGroupMembershipProtocolMapperWebhookWithManager registers the webhook for GroupMembershipProtocolMapper in the manager.
func SetupGroupMembershipProtocolMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.GroupMembershipProtocolMapper{}).
		WithValidator(&GroupMembershipProtocolMapperCustomValidator{}).
		WithDefaulter(&GroupMembershipProtocolMapperCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-openid-keycloak-operator-webhippie-de-v1alpha1-groupmembershipprotocolmapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=groupmembershipprotocolmappers,verbs=create;update,versions=v1alpha1,name=mgroupmembershipprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

// GroupMembershipProtocolMapperCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind GroupMembershipProtocolMapper when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type GroupMembershipProtocolMapperCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ admission.Defaulter[*v1alpha1.GroupMembershipProtocolMapper] = &GroupMembershipProtocolMapperCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind GroupMembershipProtocolMapper.
func (d *GroupMembershipProtocolMapperCustomDefaulter) Default(_ context.Context, mapper *v1alpha1.GroupMembershipProtocolMapper) error {
	groupmembershipprotocolmapperlog.Info("Defaulting for GroupMembershipProtocolMapper", "name", mapper.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-openid-keycloak-operator-webhippie-de-v1alpha1-groupmembershipprotocolmapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=groupmembershipprotocolmappers,verbs=create;update,versions=v1alpha1,name=vgroupmembershipprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

// GroupMembershipProtocolMapperCustomValidator struct is responsible for validating the GroupMembershipProtocolMapper resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type GroupMembershipProtocolMapperCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ admission.Validator[*v1alpha1.GroupMembershipProtocolMapper] = &GroupMembershipProtocolMapperCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type GroupMembershipProtocolMapper.
func (v *GroupMembershipProtocolMapperCustomValidator) ValidateCreate(_ context.Context, mapper *v1alpha1.GroupMembershipProtocolMapper) (admission.Warnings, error) {
	groupmembershipprotocolmapperlog.Info("Validation for GroupMembershipProtocolMapper upon creation", "name", mapper.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type GroupMembershipProtocolMapper.
func (v *GroupMembershipProtocolMapperCustomValidator) ValidateUpdate(_ context.Context, oldMapper, mapper *v1alpha1.GroupMembershipProtocolMapper) (admission.Warnings, error) {
	_ = oldMapper
	groupmembershipprotocolmapperlog.Info("Validation for GroupMembershipProtocolMapper upon update", "name", mapper.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type GroupMembershipProtocolMapper.
func (v *GroupMembershipProtocolMapperCustomValidator) ValidateDelete(_ context.Context, mapper *v1alpha1.GroupMembershipProtocolMapper) (admission.Warnings, error) {
	groupmembershipprotocolmapperlog.Info("Validation for GroupMembershipProtocolMapper upon deletion", "name", mapper.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
