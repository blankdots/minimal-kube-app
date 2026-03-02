# Tilt + kind local dev: build image, load into kind, deploy with Helm.
# Prereqs: kind cluster (make kind-create), CNPG operator, helm deps (make helm-deps once).

# Build the app image (Tilt tags it with a content-based tag for deployments)
docker_build(
  'blankdots/minimal-kube-app:latest',
  '.',
  dockerfile='Dockerfile',
)

# Load into kind the exact ref Tilt deploys (content-based tag).
# Trigger load-kind manually once after first build, or when deps change.
local_resource(
  'load-kind',
  cmd='ref=$(tilt dump image-deploy-ref blankdots/minimal-kube-app:latest 2>/dev/null) && [ -n "$ref" ] && kind load docker-image "$ref" --name minimal-kube-app || true',
  deps=['Dockerfile', 'go.mod', 'cmd/', 'internal/'],
)

# Deploy via Helm; chart includes API, cronjob, and CNPG cluster
k8s_yaml(
  helm(
    'charts/minkube',
    name='kube-app',
    namespace='default',
    values=['dev/values-tilt.yaml'],
    set=[
      'image.repository=blankdots/minimal-kube-app',
      'image.tag=latest',
      'image.pullPolicy=Never',
    ],
  ),
)

# Port-forward API so we can curl from host
k8s_resource(
  'kube-app-api',
  port_forwards=['5005:5005'],
)

# Show recent CNPG Postgres logs in Tilt UI (trigger to refresh). Uses pod label from CNPG operator.
# StatefulSet appears only after the operator reconciles the Cluster CR; pods have label cnpg.io/cluster=kube-app-cnpg
local_resource(
  'cnpg-db-logs',
  cmd='kubectl logs -l cnpg.io/cluster=kube-app-cnpg -n default --all-containers=true --tail=100 --prefix=true 2>/dev/null || echo "No CNPG pods yet (install operator and wait for cluster to be ready)."',
  resource_deps=['kube-app-api'],
  auto_init=False,
)
