help:
	@echo 'Welcome!'
	@echo ''
	@echo 'Hopefully a helpful file.'

# check go version
bootstrap: go-version-check
		GO111MODULE=off go get golang.org/x/tools/cmd/goimports

# build container
build:
	@docker build -t blankdots/minimal-kube-app .

# --- kind + Tilt (local dev) ---
# Create kind cluster (run once)
kind-create:
	@kind create cluster --config dev/kind.yaml

# Delete kind cluster
kind-delete:
	@kind delete cluster --name minimal-kube-app

# Fetch Helm chart dependencies (run once before first tilt up)
helm-deps:
	@helm repo add cnpg https://cloudnative-pg.github.io/charts
	@helm dependency build charts/minkube

# Start Tilt dev loop (build, load into kind, deploy). Requires: kind-create, helm-deps
tilt-up:
	@tilt up

# Stop Tilt
tilt-down:
	@tilt down

# One-shot: create kind cluster and install Helm deps (then run tilt up)
dev-bootstrap: kind-create helm-deps
	@echo "Run: make tilt-up"

# check go version to match one in dockerfile
go-version-check: SHELL:=/bin/bash
go-version-check:
	@GO_VERSION_MIN=$$(grep GOLANG_VERSION $(CURDIR)/Dockerfile | cut -d '-' -f2 | tr -d '}'); \
	GO_VERSION=$$(go version | grep -o 'go[0-9]\+\.[0-9]\+\(\.[0-9]\+\)\?' | tr -d 'go'); \
	IFS="." read -r -a GO_VERSION_ARR <<< "$${GO_VERSION}"; \
	IFS="." read -r -a GO_VERSION_REQ <<< "$${GO_VERSION_MIN}"; \
	if [[ $${GO_VERSION_ARR[0]} -lt $${GO_VERSION_REQ[0]} ||\
		( $${GO_VERSION_ARR[0]} -eq $${GO_VERSION_REQ[0]} &&\
		( $${GO_VERSION_ARR[1]} -lt $${GO_VERSION_REQ[1]} ||\
		( $${GO_VERSION_ARR[1]} -eq $${GO_VERSION_REQ[1]} && $${GO_VERSION_ARR[2]} -lt $${GO_VERSION_REQ[2]} )))\
	]]; then\
		echo "this requires go $${GO_VERSION_MIN} to build; found $${GO_VERSION}.";\
		exit 1;\
	fi;

# install golanci-lint and provide
# handy way to lint code
lint:
	@if ! command -v golangci-lint >/dev/null; then \
		echo "Golangci-lint needs to be installed."; \
		exit 1; \
	fi
	@echo 'Running golangci-lint'
	@golangci-lint run -E bodyclose,gocritic,gofmt,gosec,govet,nestif,nlreturn,revive,rowserrcheck

# run unit tests
test: 
	@go test -v -coverprofile=coverage.txt -covermode=atomic ./...

helm:
	@helm repo add cnpg https://cloudnative-pg.github.io/charts
	@helm dependency build charts/minkube
	@helm install kube-app charts/minkube