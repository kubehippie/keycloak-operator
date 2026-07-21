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
	"sync"
	"time"

	"github.com/kubehippie/keycloak-operator/api/common"
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// FailedKeycloakConnectionRetryPeriod is the initial requeue delay after a
	// failed connection attempt. Subsequent consecutive failures double this
	// delay (up to maxFailedKeycloakConnectionRetryPeriod) so a persistent
	// outage or upstream ban (e.g. fail2ban/WAF) does not get continuously
	// re-triggered by a tight retry loop.
	FailedKeycloakConnectionRetryPeriod    = time.Second * 10
	maxFailedKeycloakConnectionRetryPeriod = time.Minute * 10
	successKeycloakConnectionRetryPeriod   = time.Minute * 30
)

// KeycloakReconciler reconciles a Keycloak object
type KeycloakReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	failuresMu sync.Mutex
	failures   map[types.NamespacedName]int
}

// nextFailedRetryPeriod computes an exponential backoff delay for consecutive
// connection failures on the given Keycloak instance, capped at
// maxFailedKeycloakConnectionRetryPeriod.
func (r *KeycloakReconciler) nextFailedRetryPeriod(key types.NamespacedName) time.Duration {
	r.failuresMu.Lock()
	defer r.failuresMu.Unlock()

	if r.failures == nil {
		r.failures = make(map[types.NamespacedName]int)
	}

	count := r.failures[key]
	r.failures[key] = count + 1

	period := FailedKeycloakConnectionRetryPeriod << count
	if period <= 0 || period > maxFailedKeycloakConnectionRetryPeriod {
		period = maxFailedKeycloakConnectionRetryPeriod
	}

	return period
}

// resetFailedRetryPeriod clears the consecutive failure counter once a
// connection succeeds again.
func (r *KeycloakReconciler) resetFailedRetryPeriod(key types.NamespacedName) {
	r.failuresMu.Lock()
	defer r.failuresMu.Unlock()

	delete(r.failures, key)
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
		retryAfter := r.nextFailedRetryPeriod(req.NamespacedName)
		log.Info("Not connected, will retry", "retryAfter", retryAfter)
		return ctrl.Result{RequeueAfter: retryAfter}, nil
	}

	r.resetFailedRetryPeriod(req.NamespacedName)
	log.Info("Reconciling has been finished")

	return ctrl.Result{
		RequeueAfter: successKeycloakConnectionRetryPeriod,
	}, nil
}

const (
	keycloakSecretIndexField    = ".spec.secretRefs"
	keycloakConfigMapIndexField = ".spec.configMapRefs"
)

// keycloakSecretRefKeys returns the RefIndexKey values for every Secret the
// given Keycloak instance may reference (admin credentials and caCert).
func keycloakSecretRefKeys(obj client.Object) []string {
	kc, ok := obj.(*v1alpha1.Keycloak)
	if !ok {
		return nil
	}

	var keys []string
	appendSecretRef := func(ref *common.SecretKeyRefOrVal) {
		if ref == nil || ref.SecretKeyRef == nil {
			return
		}
		ns := ref.SecretKeyRef.Namespace
		if ns == "" {
			ns = kc.Namespace
		}
		keys = append(keys, RefIndexKey(ns, ref.SecretKeyRef.Name))
	}

	appendSecretRef(kc.Spec.Username)
	appendSecretRef(kc.Spec.Password)
	appendSecretRef(kc.Spec.Client)
	appendSecretRef(kc.Spec.Secret)

	if kc.Spec.CACert != nil && kc.Spec.CACert.SecretKeyRef != nil {
		ns := kc.Spec.CACert.SecretKeyRef.Namespace
		if ns == "" {
			ns = kc.Namespace
		}
		keys = append(keys, RefIndexKey(ns, kc.Spec.CACert.SecretKeyRef.Name))
	}

	return keys
}

// keycloakConfigMapRefKeys returns the RefIndexKey values for every ConfigMap
// the given Keycloak instance may reference (currently only caCert).
func keycloakConfigMapRefKeys(obj client.Object) []string {
	kc, ok := obj.(*v1alpha1.Keycloak)
	if !ok || kc.Spec.CACert == nil || kc.Spec.CACert.ConfigMapKeyRef == nil {
		return nil
	}

	ns := kc.Spec.CACert.ConfigMapKeyRef.Namespace
	if ns == "" {
		ns = kc.Namespace
	}

	return []string{RefIndexKey(ns, kc.Spec.CACert.ConfigMapKeyRef.Name)}
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeycloakReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := RegisterRefIndex(mgr, &v1alpha1.Keycloak{}, keycloakSecretIndexField, keycloakSecretRefKeys); err != nil {
		return err
	}
	if err := RegisterRefIndex(mgr, &v1alpha1.Keycloak{}, keycloakConfigMapIndexField, keycloakConfigMapRefKeys); err != nil {
		return err
	}

	newList := func() client.ObjectList { return &v1alpha1.KeycloakList{} }

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Keycloak{}).
		Watches(&corev1.Secret{}, RefEventHandler(mgr.GetClient(), newList, keycloakSecretIndexField)).
		Watches(&corev1.ConfigMap{}, RefEventHandler(mgr.GetClient(), newList, keycloakConfigMapIndexField)).
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
