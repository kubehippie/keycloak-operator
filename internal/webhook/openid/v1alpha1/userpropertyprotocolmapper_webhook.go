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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	openidv1alpha1 "github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var userpropertyprotocolmapperlog = logf.Log.WithName("userpropertyprotocolmapper-resource")

// SetupUserPropertyProtocolMapperWebhookWithManager registers the webhook for UserPropertyProtocolMapper in the manager.
func SetupUserPropertyProtocolMapperWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&openidv1alpha1.UserPropertyProtocolMapper{}).
		WithValidator(&UserPropertyProtocolMapperCustomValidator{}).
		WithDefaulter(&UserPropertyProtocolMapperCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-openid-keycloak-operator-webhippie-de-v1alpha1-userpropertyprotocolmapper,mutating=true,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=userpropertyprotocolmappers,verbs=create;update,versions=v1alpha1,name=muserpropertyprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

// UserPropertyProtocolMapperCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind UserPropertyProtocolMapper when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type UserPropertyProtocolMapperCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &UserPropertyProtocolMapperCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind UserPropertyProtocolMapper.
func (d *UserPropertyProtocolMapperCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	userpropertyprotocolmapper, ok := obj.(*openidv1alpha1.UserPropertyProtocolMapper)

	if !ok {
		return fmt.Errorf("expected an UserPropertyProtocolMapper object but got %T", obj)
	}
	userpropertyprotocolmapperlog.Info("Defaulting for UserPropertyProtocolMapper", "name", userpropertyprotocolmapper.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-openid-keycloak-operator-webhippie-de-v1alpha1-userpropertyprotocolmapper,mutating=false,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=userpropertyprotocolmappers,verbs=create;update,versions=v1alpha1,name=vuserpropertyprotocolmapper-v1alpha1.kb.io,admissionReviewVersions=v1

// UserPropertyProtocolMapperCustomValidator struct is responsible for validating the UserPropertyProtocolMapper resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type UserPropertyProtocolMapperCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &UserPropertyProtocolMapperCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type UserPropertyProtocolMapper.
func (v *UserPropertyProtocolMapperCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	userpropertyprotocolmapper, ok := obj.(*openidv1alpha1.UserPropertyProtocolMapper)
	if !ok {
		return nil, fmt.Errorf("expected a UserPropertyProtocolMapper object but got %T", obj)
	}
	userpropertyprotocolmapperlog.Info("Validation for UserPropertyProtocolMapper upon creation", "name", userpropertyprotocolmapper.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type UserPropertyProtocolMapper.
func (v *UserPropertyProtocolMapperCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	userpropertyprotocolmapper, ok := newObj.(*openidv1alpha1.UserPropertyProtocolMapper)
	if !ok {
		return nil, fmt.Errorf("expected a UserPropertyProtocolMapper object for the newObj but got %T", newObj)
	}
	userpropertyprotocolmapperlog.Info("Validation for UserPropertyProtocolMapper upon update", "name", userpropertyprotocolmapper.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type UserPropertyProtocolMapper.
func (v *UserPropertyProtocolMapperCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	userpropertyprotocolmapper, ok := obj.(*openidv1alpha1.UserPropertyProtocolMapper)
	if !ok {
		return nil, fmt.Errorf("expected a UserPropertyProtocolMapper object but got %T", obj)
	}
	userpropertyprotocolmapperlog.Info("Validation for UserPropertyProtocolMapper upon deletion", "name", userpropertyprotocolmapper.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
