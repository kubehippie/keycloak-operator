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
var openidclientlog = logf.Log.WithName("openidclient-resource")

// SetupOpenIDClientWebhookWithManager registers the webhook for OpenIDClient in the manager.
func SetupOpenIDClientWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &v1alpha1.OpenIDClient{}).
		WithValidator(&OpenIDClientCustomValidator{}).
		WithDefaulter(&OpenIDClientCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-openid-keycloak-operator-webhippie-de-v1alpha1-openidclient,mutating=true,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=openidclients,verbs=create;update,versions=v1alpha1,name=mopenidclient-v1alpha1.kb.io,admissionReviewVersions=v1

// OpenIDClientCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind OpenIDClient when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type OpenIDClientCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ admission.Defaulter[*v1alpha1.OpenIDClient] = &OpenIDClientCustomDefaulter{}

// Default implements admission.Defaulter so a webhook will be registered for the Kind OpenIDClient.
func (d *OpenIDClientCustomDefaulter) Default(_ context.Context, client *v1alpha1.OpenIDClient) error {
	openidclientlog.Info("Defaulting for OpenIDClient", "name", client.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-openid-keycloak-operator-webhippie-de-v1alpha1-openidclient,mutating=false,failurePolicy=fail,sideEffects=None,groups=openid.keycloak-operator.webhippie.de,resources=openidclients,verbs=create;update,versions=v1alpha1,name=vopenidclient-v1alpha1.kb.io,admissionReviewVersions=v1

// OpenIDClientCustomValidator struct is responsible for validating the OpenIDClient resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type OpenIDClientCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ admission.Validator[*v1alpha1.OpenIDClient] = &OpenIDClientCustomValidator{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type OpenIDClient.
func (v *OpenIDClientCustomValidator) ValidateCreate(_ context.Context, client *v1alpha1.OpenIDClient) (admission.Warnings, error) {
	openidclientlog.Info("Validation for OpenIDClient upon creation", "name", client.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type OpenIDClient.
func (v *OpenIDClientCustomValidator) ValidateUpdate(_ context.Context, oldClient, client *v1alpha1.OpenIDClient) (admission.Warnings, error) {
	_ = oldClient
	openidclientlog.Info("Validation for OpenIDClient upon update", "name", client.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type OpenIDClient.
func (v *OpenIDClientCustomValidator) ValidateDelete(_ context.Context, client *v1alpha1.OpenIDClient) (admission.Warnings, error) {
	openidclientlog.Info("Validation for OpenIDClient upon deletion", "name", client.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
