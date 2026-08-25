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
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
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

const oidcAdvancedGroupMapperFinalizer = "keycloak-operator.webhippie.de/oidcadvancedgroupmapper"

// stringPair mirrors Keycloak's internal MapperTypeSerializer.StringPair
// representation used to encode MAP_TYPE mapper config properties (such as
// the "claims" property of the "oidc-advanced-group-idp-mapper") into a
// single flat config string.
type stringPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// encodeClaims serializes a claim map into the JSON array format expected by
// Keycloak's MapperTypeSerializer for MAP_TYPE config properties.
func encodeClaims(claims map[string]string) string {
	pairs := make([]stringPair, 0, len(claims))
	keys := make([]string, 0, len(claims))
	for key := range claims {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		pairs = append(pairs, stringPair{Key: key, Value: claims[key]})
	}

	encoded, err := json.Marshal(pairs)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}

// OIDCAdvancedGroupMapperReconciler reconciles a OIDCAdvancedGroupMapper object.
type OIDCAdvancedGroupMapperReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=identity.keycloak-operator.webhippie.de,resources=oidcadvancedgroupmappers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=identity.keycloak-operator.webhippie.de,resources=oidcadvancedgroupmappers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=identity.keycloak-operator.webhippie.de,resources=oidcadvancedgroupmappers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *OIDCAdvancedGroupMapperReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &identityv1alpha1.OIDCAdvancedGroupMapper{}
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

	if !controllerutil.ContainsFinalizer(instance, oidcAdvancedGroupMapperFinalizer) {
		controllerutil.AddFinalizer(instance, oidcAdvancedGroupMapperFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	return r.reconcileOIDCAdvancedGroupMapper(ctx, instance, session, alias)
}

func (r *OIDCAdvancedGroupMapperReconciler) handleDeletion(ctx context.Context, instance *identityv1alpha1.OIDCAdvancedGroupMapper, session *controller.KeycloakSession, alias string) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, oidcAdvancedGroupMapperFinalizer) {
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

	controllerutil.RemoveFinalizer(instance, oidcAdvancedGroupMapperFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *OIDCAdvancedGroupMapperReconciler) reconcileOIDCAdvancedGroupMapper(ctx context.Context, instance *identityv1alpha1.OIDCAdvancedGroupMapper, session *controller.KeycloakSession, alias string) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	desired := oidcAdvancedGroupMapperToGocloak(instance, alias)

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

func oidcAdvancedGroupMapperToGocloak(mapper *identityv1alpha1.OIDCAdvancedGroupMapper, alias string) gocloak.IdentityProviderMapper {
	config := map[string]string{
		"claims": encodeClaims(mapper.Spec.Claims),
		"group":  mapper.Spec.Group,
	}

	if mapper.Spec.ClaimValuesRegex != nil {
		config["are.claim.values.regex"] = strconv.FormatBool(*mapper.Spec.ClaimValuesRegex)
	}

	return gocloak.IdentityProviderMapper{
		Name:                   gocloak.StringP(mapper.Spec.Name),
		IdentityProviderMapper: gocloak.StringP("oidc-advanced-group-idp-mapper"),
		IdentityProviderAlias:  gocloak.StringP(alias),
		Config:                 config,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *OIDCAdvancedGroupMapperReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&identityv1alpha1.OIDCAdvancedGroupMapper{}).
		Named("identity-oidcadvancedgroupmapper").
		Complete(r)
}
