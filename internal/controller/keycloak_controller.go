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

	"github.com/kubehippie/keycloak-operator/api/common"
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	FailedKeycloakConnectionRetryPeriod  = time.Second * 10
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
func (r *KeycloakReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &v1alpha1.Keycloak{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Instance not found")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, fmt.Errorf("unable to fetch: %w", err)
	}

	if err := r.updateConnectionStatus(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	if !instance.Status.Connected {
		log.Info("Not connected, will retry")
		return ctrl.Result{RequeueAfter: FailedKeycloakConnectionRetryPeriod}, nil
	}

	log.Info("Reconciling has been finished")

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
	log.Info("Start updating connection status")
	connected := false

	_, err := KeycloakSessionForKeycloak(
		ctx,
		r.Client,
		&common.KeycloakRef{
			Kind: "Keycloak",
			Name: instance.Name,
		},
		instance.Namespace,
	)

	if err != nil {
		log.Error(err, "Unable to connect to Keycloak")
	} else {
		connected = true
	}

	if instance.Status.Connected == connected {
		log.Info("Connection status unchanged", "status", instance.Status.Connected)
		return nil
	}

	log.Info("Connection status changed", "from", instance.Status.Connected, "to", connected)
	instance.Status.Connected = connected

	if err := r.Status().Update(ctx, instance); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	log.Info("Status have been updated", "status", instance.Status)
	return nil
}
