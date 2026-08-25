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

package identity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v14"
	identityv1alpha1 "github.com/kubehippie/keycloak-operator/api/identity/v1alpha1"
	controller "github.com/kubehippie/keycloak-operator/internal/controller"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const customProviderMapperFinalizer = "keycloak-operator.webhippie.de/customprovidermapper"

// CustomProviderMapperReconciler reconciles a CustomProviderMapper object
type CustomProviderMapperReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=identity.keycloak-operator.webhippie.de,resources=customprovidermappers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=identity.keycloak-operator.webhippie.de,resources=customprovidermappers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=identity.keycloak-operator.webhippie.de,resources=customprovidermappers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *CustomProviderMapperReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &identityv1alpha1.CustomProviderMapper{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("unable to fetch: %w", err)
	}

	session, alias, err := controller.KeycloakSessionForIdentityProvider(ctx, r.Client, instance.Spec.IdentityProviderRef, req.Namespace)
	if err != nil {
		if !instance.DeletionTimestamp.IsZero() && apierrors.IsNotFound(err) {
			log.Info("Identity provider no longer exists, skipping Keycloak cleanup", "error", err.Error())
			return r.handleDeletion(ctx, instance, nil, "")
		}
		log.Error(err, "Unable to get Keycloak session")
		return ctrl.Result{RequeueAfter: controller.FailedKeycloakConnectionRetryPeriod}, nil
	}

	if !instance.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, instance, session, alias)
	}

	if !controllerutil.ContainsFinalizer(instance, customProviderMapperFinalizer) {
		controllerutil.AddFinalizer(instance, customProviderMapperFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	return r.reconcileCustomProviderMapper(ctx, instance, session, alias)
}

func (r *CustomProviderMapperReconciler) handleDeletion(ctx context.Context, instance *identityv1alpha1.CustomProviderMapper, session *controller.KeycloakSession, alias string) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, customProviderMapperFinalizer) {
		return ctrl.Result{}, nil
	}

	if session != nil && instance.Status.KeycloakID != nil {
		log.Info("Deleting identity provider mapper from Keycloak", "id", *instance.Status.KeycloakID)
		if err := session.Client.DeleteIdentityProviderMapper(ctx, session.Token.AccessToken, session.RealmName, alias, *instance.Status.KeycloakID); err != nil {
			var apiErr *gocloak.APIError
			if errors.As(err, &apiErr) && apiErr.Code == 404 {
				log.Info("Identity provider mapper already absent in Keycloak, skipping delete")
			} else {
				return ctrl.Result{}, fmt.Errorf("failed to delete identity provider mapper from Keycloak: %w", err)
			}
		}
	}

	controllerutil.RemoveFinalizer(instance, customProviderMapperFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *CustomProviderMapperReconciler) reconcileCustomProviderMapper(ctx context.Context, instance *identityv1alpha1.CustomProviderMapper, session *controller.KeycloakSession, alias string) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	desired := customProviderMapperToGocloak(instance, alias)

	if instance.Status.KeycloakID == nil {
		existing, err := session.Client.GetIdentityProviderMappers(ctx, session.Token.AccessToken, session.RealmName, alias)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to search for existing identity provider mapper: %w", err)
		}

		var match *gocloak.IdentityProviderMapper
		for _, mapper := range existing {
			if mapper != nil && mapper.Name != nil && *mapper.Name == instance.Spec.Name {
				match = mapper
				break
			}
		}

		if match != nil {
			if match.ID == nil {
				return ctrl.Result{}, fmt.Errorf("found existing identity provider mapper %q without ID", instance.Spec.Name)
			}
			log.Info("Adopting existing Keycloak identity provider mapper", "id", *match.ID)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, match.ID); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			id, err := session.Client.CreateIdentityProviderMapper(ctx, session.Token.AccessToken, session.RealmName, alias, desired)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to create identity provider mapper in Keycloak: %w", err)
			}
			log.Info("Identity provider mapper created in Keycloak", "id", id)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, &id); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	desired.ID = instance.Status.KeycloakID
	if err := session.Client.UpdateIdentityProviderMapper(ctx, session.Token.AccessToken, session.RealmName, alias, desired); err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			log.Info("Identity provider mapper missing in Keycloak, clearing status to recreate", "id", *instance.Status.KeycloakID)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, nil); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{RequeueAfter: time.Second}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to update identity provider mapper in Keycloak: %w", err)
	}

	log.Info("Identity provider mapper reconciled", "id", *instance.Status.KeycloakID)
	return ctrl.Result{}, nil
}

func customProviderMapperToGocloak(mapper *identityv1alpha1.CustomProviderMapper, alias string) gocloak.IdentityProviderMapper {
	return gocloak.IdentityProviderMapper{
		Name:                   gocloak.StringP(mapper.Spec.Name),
		IdentityProviderMapper: gocloak.StringP(mapper.Spec.MapperType),
		IdentityProviderAlias:  gocloak.StringP(alias),
		Config:                 mapper.Spec.Config,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomProviderMapperReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&identityv1alpha1.CustomProviderMapper{}).
		Named("identity-customprovidermapper").
		Complete(r)
}
