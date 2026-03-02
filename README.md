## kube-app

### Local Development

The Makefile is primarily designed to be an aid during development work.

Start setup:
```bash
make bootstrap
```

#### kind + Tilt (recommended)

Local Kubernetes in Docker with live reload. Requires [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation), [Tilt](https://docs.tilt.dev/install.html), and the [CloudNative PG](https://cloudnative-pg.io/) operator (for the database).

```bash
# One-time: create cluster, install CNPG operator, fetch Helm chart dependencies
make kind-create
helm repo add cnpg https://cloudnative-pg.github.io/charts
helm install cnpg cnpg/cloudnative-pg -n cnpg-system --create-namespace
make helm-deps

# Start Tilt (builds image, loads into kind, deploys via Helm; port-forwards API to 5005)
make tilt-up
```

Then open the Tilt UI (default http://localhost:10350). Once the API is ready:

```bash
curl -v -H "Authorization: Bearer test" "http://localhost:5005/query?package=express"
# or: curl -v -H "X-API-Key: test" "http://localhost:5005/query?package=express"
```

Stop:

```bash
make tilt-down
make kind-delete   # when you want to remove the cluster
```

#### Build container only

```bash
make build
```

#### Unit Testing

Using [golangci-lint](https://golangci-lint.run/usage/install/) and Go tests:

```bash
make lint
make test
```

### K8s Deployment (production or manual)

Requires a Kubernetes cluster (e.g. [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)), `helm`, and `kubectl`.

#### Install and test Helm chart

```bash
make helm-deps
make helm
kubectl port-forward deployment/kube-app-api 5005:5005
curl -v -H "Authorization: Bearer test" "http://localhost:5005/query?package=express"
# or: curl -v -H "X-API-Key: test" "http://localhost:5005/query?package=express"
```
