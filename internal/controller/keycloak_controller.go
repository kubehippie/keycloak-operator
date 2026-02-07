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

package controller

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
)

const (
	failedKeycloakConnectionRetryPeriod  = time.Second * 10
	successKeycloakConnectionRetryPeriod = time.Minute * 30
)

// KeycloakReconciler reconciles a Keycloak object
type KeycloakReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=keycloaks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=keycloaks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=keycloaks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Keycloak object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.1/pkg/reconcile
func (r *KeycloakReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling Keycloak")

	instance := &v1alpha1.Keycloak{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Instance not found")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, fmt.Errorf("unable to get Keycloak: %w", err)
	}

	if err := r.updateConnectionStatus(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	if !instance.Status.Connected {
		log.Info("Keycloak is not connected, will retry")
		return ctrl.Result{RequeueAfter: failedKeycloakConnectionRetryPeriod}, nil
	}

	log.Info("Reconciling Keycloak has been finished")

	return ctrl.Result{
		RequeueAfter: successKeycloakConnectionRetryPeriod,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeycloakReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Keycloak{}).
		Named("keycloak").
		Complete(r)
}

func (r *KeycloakReconciler) updateConnectionStatus(ctx context.Context, instance *v1alpha1.Keycloak) error {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Start updating connection status to Keycloak")

	if err := r.Status().Update(ctx, instance); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	log.Info("Status have been updated", "status", instance.Status)
	return nil
}
