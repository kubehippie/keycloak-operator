allow_k8s_contexts('kind-keycloak-operator')

local_resource(
  'keycloak',
  cmd='helm upgrade --install keycloak oci://ghcr.io/codecentric/helm-charts/keycloakx --create-namespace --namespace keycloak-system --values test/e2e/testdata/keycloak.yaml --set image.tag=26.6 --timeout=5m --wait --hide-notes',
  deps=['test/e2e/testdata/keycloak.yaml'],
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

local_resource(
  'generate',
  'make generate',
  deps=[
    'api/**/**/*_types.go',
    'api/**/groupversion_info.go',
    'internal/',
    'cmd/',
    'hack/boilerplate.go.txt',
  ],
)

local_resource(
  'manifests',
  'make manifests',
  deps=[
    'api/**/**/*_types.go',
    'api/**/groupversion_info.go',
    'internal/',
  ],
  resource_deps=[
    'generate',
  ],
)

k8s_yaml(
  kustomize('config/default')
)

k8s_resource(
  'keycloak-operator-controller-manager',
  extra_pod_selectors=[{'control-plane': 'controller-manager'}],
  resource_deps=['keycloak'],
)
