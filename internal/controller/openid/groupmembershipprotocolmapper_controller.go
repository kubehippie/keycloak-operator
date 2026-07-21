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

package openid

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v14"
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/openid/v1alpha1"
	controller "github.com/kubehippie/keycloak-operator/internal/controller"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const groupMembershipProtocolMapperFinalizer = "keycloak-operator.webhippie.de/groupmembershipprotocolmapper"

// GroupMembershipProtocolMapperReconciler reconciles a GroupMembershipProtocolMapper object.
type GroupMembershipProtocolMapperReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=groupmembershipprotocolmappers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=groupmembershipprotocolmappers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=openid.keycloak-operator.webhippie.de,resources=groupmembershipprotocolmappers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *GroupMembershipProtocolMapperReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &v1alpha1.GroupMembershipProtocolMapper{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("unable to fetch: %w", err)
	}

	session, idOfClient, err := controller.KeycloakSessionForClient(ctx, r.Client, instance.Spec.ClientRef, req.Namespace)
	if err != nil {
		log.Error(err, "Unable to get Keycloak session")
		return ctrl.Result{RequeueAfter: controller.FailedKeycloakConnectionRetryPeriod}, nil
	}

	if !instance.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, instance, session, idOfClient)
	}

	if !controllerutil.ContainsFinalizer(instance, groupMembershipProtocolMapperFinalizer) {
		controllerutil.AddFinalizer(instance, groupMembershipProtocolMapperFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return r.reconcileGroupMembershipProtocolMapper(ctx, instance, session, idOfClient)
}

func (r *GroupMembershipProtocolMapperReconciler) handleDeletion(ctx context.Context, instance *v1alpha1.GroupMembershipProtocolMapper, session *controller.KeycloakSession, idOfClient string) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, groupMembershipProtocolMapperFinalizer) {
		return ctrl.Result{}, nil
	}

	if instance.Status.KeycloakID != nil {
		log.Info("Deleting protocol mapper from Keycloak", "id", *instance.Status.KeycloakID)
		if err := session.Client.DeleteClientProtocolMapper(ctx, session.Token.AccessToken, session.RealmName, idOfClient, *instance.Status.KeycloakID); err != nil {
			if !isNotFoundAPIError(err) {
				return ctrl.Result{}, fmt.Errorf("failed to delete protocol mapper from Keycloak: %w", err)
			}
			log.Info("Protocol mapper already absent in Keycloak, skipping delete")
		}
	}

	controllerutil.RemoveFinalizer(instance, groupMembershipProtocolMapperFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *GroupMembershipProtocolMapperReconciler) reconcileGroupMembershipProtocolMapper(ctx context.Context, instance *v1alpha1.GroupMembershipProtocolMapper, session *controller.KeycloakSession, idOfClient string) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	desired := groupMembershipProtocolMapperToGocloak(instance)

	if instance.Status.KeycloakID == nil {
		existingClient, err := session.Client.GetClient(ctx, session.Token.AccessToken, session.RealmName, idOfClient)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to search for existing protocol mapper: %w", err)
		}

		if match := findClientProtocolMapperByName(existingClient, instance.Spec.Name); match != nil {
			log.Info("Adopting existing Keycloak protocol mapper", "id", *match.ID)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, match.ID); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			id, err := session.Client.CreateClientProtocolMapper(ctx, session.Token.AccessToken, session.RealmName, idOfClient, desired)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to create protocol mapper in Keycloak: %w", err)
			}
			log.Info("Protocol mapper created in Keycloak", "id", id)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, &id); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	desired.ID = instance.Status.KeycloakID
	if err := session.Client.UpdateClientProtocolMapper(ctx, session.Token.AccessToken, session.RealmName, idOfClient, *instance.Status.KeycloakID, desired); err != nil {
		if isNotFoundAPIError(err) {
			log.Info("Protocol mapper missing in Keycloak, clearing status to recreate", "id", *instance.Status.KeycloakID)
			if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, nil); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to update protocol mapper in Keycloak: %w", err)
	}

	log.Info("Protocol mapper reconciled", "id", *instance.Status.KeycloakID)
	return ctrl.Result{}, nil
}

func groupMembershipProtocolMapperToGocloak(mapper *v1alpha1.GroupMembershipProtocolMapper) gocloak.ProtocolMapperRepresentation {
	return gocloak.ProtocolMapperRepresentation{
		Name:           gocloak.StringP(mapper.Spec.Name),
		Protocol:       gocloak.StringP(openidConnectProtocol),
		ProtocolMapper: gocloak.StringP("oidc-group-membership-mapper"),
		Config: map[string]string{
			"full.path":                      boolStringDefaultTrue(mapper.Spec.FullPath),
			configKeyIDTokenClaim:            boolStringDefaultTrue(mapper.Spec.AddToIDToken),
			configKeyAccessTokenClaim:        boolStringDefaultTrue(mapper.Spec.AddToAccessToken),
			configKeyUserinfoTokenClaim:      boolStringDefaultTrue(mapper.Spec.AddToUserInfo),
			configKeyIntrospectionTokenClaim: boolStringDefaultTrue(mapper.Spec.AddToTokenIntrospection),
			configKeyClaimName:               mapper.Spec.ClaimName,
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *GroupMembershipProtocolMapperReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.GroupMembershipProtocolMapper{}).
		Named("openid-groupmembershipprotocolmapper").
		Complete(r)
}
