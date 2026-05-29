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

	"github.com/Nerzal/gocloak/v13"
	"github.com/kubehippie/keycloak-operator/api/common"
	v1alpha1 "github.com/kubehippie/keycloak-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// keycloakIDObject is implemented by any CR whose status exposes the
// Keycloak-assigned UUID so it can be read and written generically.
type keycloakIDObject interface {
	client.Object
	GetKeycloakID() *string
	SetKeycloakID(id *string)
}

// updateKeycloakIDStatus sets status.keycloakID to newID and persists the
// change via the status subresource. When the current value already equals
// newID the API call is skipped and nil is returned, mirroring the
// no-op behaviour of updateConnectionStatus in the Keycloak controller.
func updateKeycloakIDStatus(ctx context.Context, c client.Client, obj keycloakIDObject, newID *string) error {
	current := obj.GetKeycloakID()
	if current != nil && newID != nil && *current == *newID {
		return nil
	}
	obj.SetKeycloakID(newID)
	if err := c.Status().Update(ctx, obj); err != nil {
		return fmt.Errorf("failed to update Keycloak ID status: %w", err)
	}
	return nil
}

// keycloakSession holds an authenticated GoCloak client together with the
// access token and the target realm name derived from the referenced Realm CR.
type keycloakSession struct {
	Client    *gocloak.GoCloak
	Token     *gocloak.JWT
	RealmName string
}

// keycloakSessionForRealm resolves a RealmRef to an authenticated Keycloak
// API session. It walks RealmRef → Realm → KeycloakRef → Keycloak, reads
// credentials, and returns a ready-to-use session.
// defaultNamespace is used when a ref carries no explicit namespace.
func keycloakSessionForRealm(ctx context.Context, c client.Client, realmRef *common.RealmRef, defaultNamespace string) (*keycloakSession, error) {
	realmNS := realmRef.Namespace
	if realmNS == "" {
		realmNS = defaultNamespace
	}

	realm := &v1alpha1.Realm{}
	if err := c.Get(ctx, types.NamespacedName{Namespace: realmNS, Name: realmRef.Name}, realm); err != nil {
		return nil, fmt.Errorf("unable to fetch realm %s/%s: %w", realmNS, realmRef.Name, err)
	}

	session, err := keycloakSessionForKeycloak(ctx, c, realm.Spec.KeycloakRef, realmNS)
	if err != nil {
		return nil, err
	}

	session.RealmName = realm.Spec.RealmName
	return session, nil
}

// keycloakSessionForKeycloak creates an authenticated Keycloak API session
// directly from a KeycloakRef. The RealmName field in the returned session is
// left empty; callers supply the target realm name themselves.
// defaultNamespace is used when the ref carries no explicit namespace.
func keycloakSessionForKeycloak(ctx context.Context, c client.Client, kcRef *common.KeycloakRef, defaultNamespace string) (*keycloakSession, error) {
	kcNS := kcRef.Namespace
	if kcNS == "" {
		kcNS = defaultNamespace
	}

	kc := &v1alpha1.Keycloak{}
	if err := c.Get(ctx, types.NamespacedName{Namespace: kcNS, Name: kcRef.Name}, kc); err != nil {
		return nil, fmt.Errorf("unable to fetch keycloak %s/%s: %w", kcNS, kcRef.Name, err)
	}

	username, err := resolveSecretKeyRefOrVal(ctx, c, kc.Spec.Username, kcNS)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve username: %w", err)
	}

	password, err := resolveSecretKeyRefOrVal(ctx, c, kc.Spec.Password, kcNS)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve password: %w", err)
	}

	kcClient := gocloak.NewClient(kc.Spec.URL)

	token, err := kcClient.LoginAdmin(ctx, username, password, kc.Spec.RealmName)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with keycloak: %w", err)
	}

	return &keycloakSession{
		Client: kcClient,
		Token:  token,
	}, nil
}

// resolveSecretKeyRefOrVal returns the string value from a SecretKeyRefOrVal.
// If the inline Value field is non-empty it is returned directly; otherwise
// the value is read from the referenced Kubernetes Secret.
func resolveSecretKeyRefOrVal(ctx context.Context, c client.Client, ref *common.SecretKeyRefOrVal, defaultNamespace string) (string, error) {
	if ref == nil {
		return "", fmt.Errorf("ref is nil")
	}

	if ref.Value != "" {
		return ref.Value, nil
	}

	if ref.SecretKeyRef == nil {
		return "", fmt.Errorf("either value or secretKeyRef must be set")
	}

	if ref.SecretKeyRef.Name == "" {
		return "", fmt.Errorf("secretKeyRef.name must be set")
	}

	if ref.SecretKeyRef.Key == "" {
		return "", fmt.Errorf("secretKeyRef.key must be set")
	}

	ns := ref.SecretKeyRef.Namespace
	if ns == "" {
		ns = defaultNamespace
	}

	secret := &corev1.Secret{}
	if err := c.Get(ctx, types.NamespacedName{Namespace: ns, Name: ref.SecretKeyRef.Name}, secret); err != nil {
		return "", fmt.Errorf("unable to get secret %s/%s: %w", ns, ref.SecretKeyRef.Name, err)
	}

	val, ok := secret.Data[ref.SecretKeyRef.Key]
	if !ok {
		return "", fmt.Errorf("key %q not found in secret %s/%s", ref.SecretKeyRef.Key, ns, ref.SecretKeyRef.Name)
	}

	return string(val), nil
}
