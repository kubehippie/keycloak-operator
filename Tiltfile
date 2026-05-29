load('ext://helm_resource', 'helm_resource')
allow_k8s_contexts('kind-keycloak-operator')
update_settings(k8s_upsert_timeout_secs=5*60)

helm_resource(
  'keycloak',
  'oci://ghcr.io/codecentric/helm-charts/keycloakx',
  namespace='keycloak-system',
  flags=[
    '--create-namespace',
    '--values=test/e2e/testdata/keycloak.yaml',
    '--timeout=5m',
    '--wait',
    '--hide-notes',
  ],
  labels=['dependencies'],
)

helm_resource(
  'cert-manager',
  'oci://quay.io/jetstack/charts/cert-manager',
  namespace='cert-manager',
  flags=[
    '--create-namespace',
    '--values=test/e2e/testdata/cert-manager.yaml',
    '--timeout=5m',
    '--wait',
    '--hide-notes',
  ],
  labels=['dependencies'],
)

docker_build(
  'ghcr.io/kubehippie/keycloak-operator',
  '.',
  dockerfile='Dockerfile.tilt',
  entrypoint='/manager',
  live_update=[
    sync('.', '/workspace'),
    run(
      'go build -o /workspace/manager ./cmd/main.go',
      trigger=['cmd/', 'internal/', 'api/', 'go.mod', 'go.sum']
    ),
  ],
)

k8s_yaml(kustomize('config/default'))

k8s_resource(
  objects=[
    'keycloak-operator-system:Namespace',
    'attributeimportermappers.identity.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'audienceprotocolmappers.openid.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'clientscopes.openid.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'defaultscopes.openid.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'fullnameprotocolmappers.openid.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'groupmembershipprotocolmappers.openid.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'groups.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'hardcodedrolemappers.identity.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'keycloaks.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'oidcadvancedgroupmappers.identity.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'oidcidentityproviders.identity.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'openidclients.openid.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'realms.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'userattributeprotocolmappers.openid.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'userpropertyprotocolmappers.openid.keycloak-operator.webhippie.de:CustomResourceDefinition',
    'users.keycloak-operator.webhippie.de:CustomResourceDefinition',
  ],
  new_name='operator-crds',
  labels=['operator'],
)

k8s_resource(
  objects=[
    'keycloak-operator-group-admin-role:ClusterRole:default',
    'keycloak-operator-group-editor-role:ClusterRole:default',
    'keycloak-operator-group-viewer-role:ClusterRole:default',
    'keycloak-operator-identity-attributeimportermapper-admin-role:ClusterRole:default',
    'keycloak-operator-identity-attributeimportermapper-editor-role:ClusterRole:default',
    'keycloak-operator-identity-attributeimportermapper-viewer-role:ClusterRole:default',
    'keycloak-operator-identity-hardcodedrolemapper-admin-role:ClusterRole:default',
    'keycloak-operator-identity-hardcodedrolemapper-editor-role:ClusterRole:default',
    'keycloak-operator-identity-hardcodedrolemapper-viewer-role:ClusterRole:default',
    'keycloak-operator-identity-oidcadvancedgroupmapper-admin-role:ClusterRole:default',
    'keycloak-operator-identity-oidcadvancedgroupmapper-editor-role:ClusterRole:default',
    'keycloak-operator-identity-oidcadvancedgroupmapper-viewer-role:ClusterRole:default',
    'keycloak-operator-identity-oidcidentityprovider-admin-role:ClusterRole:default',
    'keycloak-operator-identity-oidcidentityprovider-editor-role:ClusterRole:default',
    'keycloak-operator-identity-oidcidentityprovider-viewer-role:ClusterRole:default',
    'keycloak-operator-keycloak-admin-role:ClusterRole:default',
    'keycloak-operator-keycloak-editor-role:ClusterRole:default',
    'keycloak-operator-keycloak-viewer-role:ClusterRole:default',
    'keycloak-operator-manager-role:ClusterRole:default',
    'keycloak-operator-metrics-auth-role:ClusterRole:default',
    'keycloak-operator-metrics-reader:ClusterRole:default',
    'keycloak-operator-openid-audienceprotocolmapper-admin-role:ClusterRole:default',
    'keycloak-operator-openid-audienceprotocolmapper-editor-role:ClusterRole:default',
    'keycloak-operator-openid-audienceprotocolmapper-viewer-role:ClusterRole:default',
    'keycloak-operator-openid-clientscope-admin-role:ClusterRole:default',
    'keycloak-operator-openid-clientscope-editor-role:ClusterRole:default',
    'keycloak-operator-openid-clientscope-viewer-role:ClusterRole:default',
    'keycloak-operator-openid-defaultscopes-admin-role:ClusterRole:default',
    'keycloak-operator-openid-defaultscopes-editor-role:ClusterRole:default',
    'keycloak-operator-openid-defaultscopes-viewer-role:ClusterRole:default',
    'keycloak-operator-openid-fullnameprotocolmapper-admin-role:ClusterRole:default',
    'keycloak-operator-openid-fullnameprotocolmapper-editor-role:ClusterRole:default',
    'keycloak-operator-openid-fullnameprotocolmapper-viewer-role:ClusterRole:default',
    'keycloak-operator-openid-groupmembershipprotocolmapper-admin-role:ClusterRole:default',
    'keycloak-operator-openid-groupmembershipprotocolmapper-editor-role:ClusterRole:default',
    'keycloak-operator-openid-groupmembershipprotocolmapper-viewer-role:ClusterRole:default',
    'keycloak-operator-openid-openidclient-admin-role:ClusterRole:default',
    'keycloak-operator-openid-openidclient-editor-role:ClusterRole:default',
    'keycloak-operator-openid-openidclient-viewer-role:ClusterRole:default',
    'keycloak-operator-openid-userattributeprotocolmapper-admin-role:ClusterRole:default',
    'keycloak-operator-openid-userattributeprotocolmapper-editor-role:ClusterRole:default',
    'keycloak-operator-openid-userattributeprotocolmapper-viewer-role:ClusterRole:default',
    'keycloak-operator-openid-userpropertyprotocolmapper-admin-role:ClusterRole:default',
    'keycloak-operator-openid-userpropertyprotocolmapper-editor-role:ClusterRole:default',
    'keycloak-operator-openid-userpropertyprotocolmapper-viewer-role:ClusterRole:default',
    'keycloak-operator-realm-admin-role:ClusterRole:default',
    'keycloak-operator-realm-editor-role:ClusterRole:default',
    'keycloak-operator-realm-viewer-role:ClusterRole:default',
    'keycloak-operator-user-admin-role:ClusterRole:default',
    'keycloak-operator-user-editor-role:ClusterRole:default',
    'keycloak-operator-user-viewer-role:ClusterRole:default',
    'keycloak-operator-manager-rolebinding:ClusterRoleBinding:default',
    'keycloak-operator-metrics-auth-rolebinding:ClusterRoleBinding:default',
  ],
  new_name='operator-clusterroles',
  labels=['operator'],
)

k8s_resource(
  objects=[
    'keycloak-operator-selfsigned-issuer:Issuer:keycloak-operator-system',
    'keycloak-operator-metrics-certs:Certificate:keycloak-operator-system',
    'keycloak-operator-serving-cert:Certificate:keycloak-operator-system',
    'keycloak-operator-mutating-webhook-configuration:MutatingWebhookConfiguration:default',
    'keycloak-operator-validating-webhook-configuration:ValidatingWebhookConfiguration:default',
  ],
  new_name='operator-certs',
  resource_deps=['cert-manager', 'operator-crds'],
  labels=['operator'],
)

k8s_resource(
  'keycloak-operator-controller-manager',
  new_name='operator-manager',
  extra_pod_selectors=[{'control-plane': 'controller-manager'}],
  resource_deps=['keycloak', 'operator-certs'],
  labels=['operator'],
)
