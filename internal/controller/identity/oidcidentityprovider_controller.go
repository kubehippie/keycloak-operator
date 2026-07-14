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

	"github.com/Nerzal/gocloak/v14"
	"github.com/kubehippie/keycloak-operator/api/identity/v1alpha1"
	"github.com/kubehippie/keycloak-operator/internal/controller"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const oidcIdentityProviderFinalizer = "keycloak-operator.webhippie.de/oidcidentityprovider"

// OIDCIdentityProviderReconciler reconciles a OIDCIdentityProvider object
type OIDCIdentityProviderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=identity.keycloak-operator.webhippie.de,resources=oidcidentityproviders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=identity.keycloak-operator.webhippie.de,resources=oidcidentityproviders/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=identity.keycloak-operator.webhippie.de,resources=oidcidentityproviders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *OIDCIdentityProviderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling")

	instance := &v1alpha1.OIDCIdentityProvider{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("unable to fetch: %w", err)
	}

	session, err := controller.KeycloakSessionForRealm(ctx, r.Client, instance.Spec.RealmRef, req.Namespace)
	if err != nil {
		log.Error(err, "Unable to get Keycloak session")
		return ctrl.Result{RequeueAfter: controller.FailedKeycloakConnectionRetryPeriod}, nil
	}

	if !instance.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, instance, session)
	}

	if !controllerutil.ContainsFinalizer(instance, oidcIdentityProviderFinalizer) {
		controllerutil.AddFinalizer(instance, oidcIdentityProviderFinalizer)
		if err := r.Update(ctx, instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return r.reconcileOIDCIdentityProvider(ctx, instance, session)
}

func (r *OIDCIdentityProviderReconciler) handleDeletion(ctx context.Context, instance *v1alpha1.OIDCIdentityProvider, session *controller.KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(instance, oidcIdentityProviderFinalizer) {
		return ctrl.Result{}, nil
	}

	log.Info("Deleting identity provider from Keycloak", "alias", instance.Spec.Alias)
	if err := session.Client.DeleteIdentityProvider(ctx, session.Token.AccessToken, session.RealmName, instance.Spec.Alias); err != nil {
		var apiErr *gocloak.APIError
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			log.Info("Identity provider already absent in Keycloak, skipping delete")
		} else {
			return ctrl.Result{}, fmt.Errorf("failed to delete identity provider from Keycloak: %w", err)
		}
	}

	controllerutil.RemoveFinalizer(instance, oidcIdentityProviderFinalizer)
	if err := r.Update(ctx, instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *OIDCIdentityProviderReconciler) reconcileOIDCIdentityProvider(ctx context.Context, instance *v1alpha1.OIDCIdentityProvider, session *controller.KeycloakSession) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	desired, err := oidcIdentityProviderToGocloak(ctx, r.Client, instance, instance.Namespace)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to build identity provider representation: %w", err)
	}

	existing, err := session.Client.GetIdentityProvider(ctx, session.Token.AccessToken, session.RealmName, instance.Spec.Alias)
	if err != nil {
		var apiErr *gocloak.APIError
		if !errors.As(err, &apiErr) || apiErr.Code != 404 {
			return ctrl.Result{}, fmt.Errorf("failed to check for existing identity provider: %w", err)
		}

		id, err := session.Client.CreateIdentityProvider(ctx, session.Token.AccessToken, session.RealmName, desired)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create identity provider in Keycloak: %w", err)
		}
		log.Info("Identity provider created in Keycloak", "id", id)
		if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, &id); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if err := controller.UpdateKeycloakIDStatus(ctx, r.Client, instance, existing.InternalID); err != nil {
		return ctrl.Result{}, err
	}

	if err := session.Client.UpdateIdentityProvider(ctx, session.Token.AccessToken, session.RealmName, instance.Spec.Alias, desired); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update identity provider in Keycloak: %w", err)
	}

	log.Info("Identity provider reconciled", "alias", instance.Spec.Alias)
	return ctrl.Result{}, nil
}

// oidcIdentityProviderToGocloak converts an OIDCIdentityProvider CR spec into
// the gocloak representation used for create and update API calls. The
// client secret (if configured as a Secret reference) is resolved from the
// cluster at call time.
func oidcIdentityProviderToGocloak(ctx context.Context, cl client.Client, p *v1alpha1.OIDCIdentityProvider, ns string) (gocloak.IdentityProviderRepresentation, error) {
	config := map[string]string{
		"authorizationUrl": p.Spec.AuthorizationURL,
		"tokenUrl":         p.Spec.TokenURL,
		"clientId":         p.Spec.ClientID,
	}

	clientSecret, err := controller.ResolveSecretKeyRefOrVal(ctx, cl, p.Spec.ClientSecret, ns)
	if err != nil {
		return gocloak.IdentityProviderRepresentation{}, fmt.Errorf("unable to resolve clientSecret: %w", err)
	}
	config["clientSecret"] = clientSecret

	if p.Spec.UserInfoURL != nil {
		config["userInfoUrl"] = *p.Spec.UserInfoURL
	}
	if p.Spec.JwksURL != nil {
		config["jwksUrl"] = *p.Spec.JwksURL
		config["useJwksUrl"] = "true"
	}
	if p.Spec.LogoutURL != nil {
		config["logoutUrl"] = *p.Spec.LogoutURL
	}
	if p.Spec.Issuer != nil {
		config["issuer"] = *p.Spec.Issuer
	}
	if p.Spec.DefaultScopes != nil {
		config["defaultScope"] = *p.Spec.DefaultScopes
	}
	if p.Spec.ValidateSignature != nil {
		config["validateSignature"] = boolToString(*p.Spec.ValidateSignature)
	}
	if p.Spec.BackchannelSupported != nil {
		config["backchannelSupported"] = boolToString(*p.Spec.BackchannelSupported)
	}
	if p.Spec.SyncMode != nil {
		config["syncMode"] = *p.Spec.SyncMode
	}

	idp := gocloak.IdentityProviderRepresentation{
		Alias:                     gocloak.StringP(p.Spec.Alias),
		ProviderID:                gocloak.StringP("oidc"),
		Enabled:                   p.Spec.Enabled,
		DisplayName:               p.Spec.DisplayName,
		HideOnLogin:               p.Spec.HideOnLoginPage,
		StoreToken:                p.Spec.StoreToken,
		TrustEmail:                p.Spec.TrustEmail,
		LinkOnly:                  p.Spec.LinkOnly,
		FirstBrokerLoginFlowAlias: p.Spec.FirstBrokerLoginFlowAlias,
		PostBrokerLoginFlowAlias:  p.Spec.PostBrokerLoginFlowAlias,
		Config:                    config,
	}

	return idp, nil
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// SetupWithManager sets up the controller with the Manager.
func (r *OIDCIdentityProviderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.OIDCIdentityProvider{}).
		Named("identity-oidcidentityprovider").
		Complete(r)
}
