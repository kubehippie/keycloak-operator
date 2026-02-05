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

	keycloakoperatorwebhippiedev1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var clusterkeycloaklog = logf.Log.WithName("clusterkeycloak-resource")

// SetupClusterKeycloakWebhookWithManager registers the webhook for ClusterKeycloak in the manager.
func SetupClusterKeycloakWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&keycloakoperatorwebhippiedev1alpha1.ClusterKeycloak{}).
		WithValidator(&ClusterKeycloakCustomValidator{}).
		WithDefaulter(&ClusterKeycloakCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-keycloak-operator-webhippie-de-v1alpha1-clusterkeycloak,mutating=true,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=clusterkeycloaks,verbs=create;update,versions=v1alpha1,name=mclusterkeycloak-v1alpha1.kb.io,admissionReviewVersions=v1

// ClusterKeycloakCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind ClusterKeycloak when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type ClusterKeycloakCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &ClusterKeycloakCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind ClusterKeycloak.
func (d *ClusterKeycloakCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	clusterkeycloak, ok := obj.(*keycloakoperatorwebhippiedev1alpha1.ClusterKeycloak)

	if !ok {
		return fmt.Errorf("expected an ClusterKeycloak object but got %T", obj)
	}
	clusterkeycloaklog.Info("Defaulting for ClusterKeycloak", "name", clusterkeycloak.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-keycloak-operator-webhippie-de-v1alpha1-clusterkeycloak,mutating=false,failurePolicy=fail,sideEffects=None,groups=keycloak-operator.webhippie.de,resources=clusterkeycloaks,verbs=create;update,versions=v1alpha1,name=vclusterkeycloak-v1alpha1.kb.io,admissionReviewVersions=v1

// ClusterKeycloakCustomValidator struct is responsible for validating the ClusterKeycloak resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ClusterKeycloakCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &ClusterKeycloakCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type ClusterKeycloak.
func (v *ClusterKeycloakCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	clusterkeycloak, ok := obj.(*keycloakoperatorwebhippiedev1alpha1.ClusterKeycloak)
	if !ok {
		return nil, fmt.Errorf("expected a ClusterKeycloak object but got %T", obj)
	}
	clusterkeycloaklog.Info("Validation for ClusterKeycloak upon creation", "name", clusterkeycloak.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type ClusterKeycloak.
func (v *ClusterKeycloakCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	clusterkeycloak, ok := newObj.(*keycloakoperatorwebhippiedev1alpha1.ClusterKeycloak)
	if !ok {
		return nil, fmt.Errorf("expected a ClusterKeycloak object for the newObj but got %T", newObj)
	}
	clusterkeycloaklog.Info("Validation for ClusterKeycloak upon update", "name", clusterkeycloak.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type ClusterKeycloak.
func (v *ClusterKeycloakCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	clusterkeycloak, ok := obj.(*keycloakoperatorwebhippiedev1alpha1.ClusterKeycloak)
	if !ok {
		return nil, fmt.Errorf("expected a ClusterKeycloak object but got %T", obj)
	}
	clusterkeycloaklog.Info("Validation for ClusterKeycloak upon deletion", "name", clusterkeycloak.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
