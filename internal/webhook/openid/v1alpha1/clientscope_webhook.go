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
var clientscopelog = logf.Log.WithName("clientscope-resource")

// SetupClientScopeWebhookWithManager registers the webhook for ClientScope in the manager.
func SetupClientScopeWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&openidv1alpha1.ClientScope{}).
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

var _ webhook.CustomDefaulter = &ClientScopeCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind ClientScope.
func (d *ClientScopeCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	clientscope, ok := obj.(*openidv1alpha1.ClientScope)

	if !ok {
		return fmt.Errorf("expected an ClientScope object but got %T", obj)
	}
	clientscopelog.Info("Defaulting for ClientScope", "name", clientscope.GetName())

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

var _ webhook.CustomValidator = &ClientScopeCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type ClientScope.
func (v *ClientScopeCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	clientscope, ok := obj.(*openidv1alpha1.ClientScope)
	if !ok {
		return nil, fmt.Errorf("expected a ClientScope object but got %T", obj)
	}
	clientscopelog.Info("Validation for ClientScope upon creation", "name", clientscope.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type ClientScope.
func (v *ClientScopeCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	clientscope, ok := newObj.(*openidv1alpha1.ClientScope)
	if !ok {
		return nil, fmt.Errorf("expected a ClientScope object for the newObj but got %T", newObj)
	}
	clientscopelog.Info("Validation for ClientScope upon update", "name", clientscope.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type ClientScope.
func (v *ClientScopeCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	clientscope, ok := obj.(*openidv1alpha1.ClientScope)
	if !ok {
		return nil, fmt.Errorf("expected a ClientScope object but got %T", obj)
	}
	clientscopelog.Info("Validation for ClientScope upon deletion", "name", clientscope.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
