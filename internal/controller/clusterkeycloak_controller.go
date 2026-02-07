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
	failedClusterKeycloakConnectionRetryPeriod  = time.Second * 10
	successClusterKeycloakConnectionRetryPeriod = time.Minute * 30
)

// ClusterKeycloakReconciler reconciles a ClusterKeycloak object
type ClusterKeycloakReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=clusterkeycloaks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=clusterkeycloaks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=keycloak-operator.webhippie.de,resources=clusterkeycloaks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterKeycloak object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.1/pkg/reconcile
func (r *ClusterKeycloakReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling ClusterKeycloak")

	instance := &v1alpha1.ClusterKeycloak{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Instance not found")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, fmt.Errorf("unable to get ClusterKeycloak: %w", err)
	}

	if err := r.updateConnectionStatus(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	if !instance.Status.Connected {
		log.Info("ClusterKeycloak is not connected, will retry")
		return ctrl.Result{RequeueAfter: failedClusterKeycloakConnectionRetryPeriod}, nil
	}

	log.Info("Reconciling ClusterKeycloak has been finished")

	return ctrl.Result{
		RequeueAfter: successClusterKeycloakConnectionRetryPeriod,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterKeycloakReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ClusterKeycloak{}).
		Named("clusterkeycloak").
		Complete(r)
}

func (r *ClusterKeycloakReconciler) updateConnectionStatus(ctx context.Context, instance *v1alpha1.ClusterKeycloak) error {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Start updating connection status to ClusterKeycloak")

	if err := r.Status().Update(ctx, instance); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	log.Info("Status have been updated", "status", instance.Status)
	return nil
}
