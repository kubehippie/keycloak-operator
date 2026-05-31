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
var clientscopelog = logf.Log.WithName("clientscope-resource")

// SetupClientScopeWebhookWithManager registers the webhook for ClientScope in the manager.
func SetupClientScopeWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.ClientScope{}).
		WithValidator(&ClientScopeCustomValidator{}).
		WithDefaulter(&ClientScopeCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-openid-keycloak-operator-webhippie-de-v1alpha1-clientscope,mutating=true,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=clientscopes,verbs=create;update,versions=v1alpha1,name=mclientscope-v1alpha1.kb.io,admissionReviewVersions=v1

// ClientScopeCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind ClientScope when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type ClientScopeCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ admission.Defaulter[*v1alpha1.ClientScope] = &ClientScopeCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind ClientScope.
func (d *ClientScopeCustomDefaulter) Default(_ context.Context, clientScope *v1alpha1.ClientScope) error {
	clientscopelog.Info("Defaulting for ClientScope", "name", clientScope.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-openid-keycloak-operator-webhippie-de-v1alpha1-clientscope,mutating=false,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=clientscopes,verbs=create;update,versions=v1alpha1,name=vclientscope-v1alpha1.kb.io,admissionReviewVersions=v1

// ClientScopeCustomValidator struct is responsible for validating the ClientScope resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ClientScopeCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ admission.Validator[*v1alpha1.ClientScope] = &ClientScopeCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type ClientScope.
func (v *ClientScopeCustomValidator) ValidateCreate(_ context.Context, clientScope *v1alpha1.ClientScope) (admission.Warnings, error) {
	clientscopelog.Info("Validation for ClientScope upon creation", "name", clientScope.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type ClientScope.
func (v *ClientScopeCustomValidator) ValidateUpdate(_ context.Context, oldClientScope, clientScope *v1alpha1.ClientScope) (admission.Warnings, error) {
	_ = oldClientScope
	clientscopelog.Info("Validation for ClientScope upon update", "name", clientScope.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type ClientScope.
func (v *ClientScopeCustomValidator) ValidateDelete(_ context.Context, clientScope *v1alpha1.ClientScope) (admission.Warnings, error) {
	clientscopelog.Info("Validation for ClientScope upon deletion", "name", clientScope.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
