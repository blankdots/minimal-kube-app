# kube-app

![Version: 0.0.2](https://img.shields.io/badge/Version-0.0.2-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.4](https://img.shields.io/badge/AppVersion-0.0.4-informational?style=flat-square)

kube-app exercise helm chart

## Prerequisites

Install the [CloudNative PG](https://cloudnative-pg.io/) operator before installing this chart (required for the database cluster):

```bash
helm repo add cnpg https://cloudnative-pg.github.io/charts
helm install cnpg cnpg/cloudnative-pg -n cnpg-system --create-namespace
```

## Source Code

* <https://github.com/blankdots/minimal-kube-app>

## Requirements

Kubernetes: `>= 1.26.0`

| Repository | Name | Version |
|------------|------|---------|
| https://cloudnative-pg.github.io/charts | cluster (alias: cnpg) | 0.5.x |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| api.token | string | `""` |  |
| cronjob.api | string | `""` |  |
| cronjob.backoffLimit | int | `1` |  |
| cronjob.concurrencyPolicy | int | `1` |  |
| cronjob.failedJobsHistoryLimit | int | `1` |  |
| cronjob.schedule | string | `""` |  |
| cronjob.successfulJobsHistoryLimit | int | `1` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"ghcr.io/blankdots/minimal-kube-app"` |  |
| image.tag | string | `""` |  |
| imagePullSecrets | list | `[]` |  |
| gateway.enabled | bool | `false` |  |
| gateway.gatewayRef.name | string | `""` | Gateway to attach HTTPRoute to |
| gateway.gatewayRef.namespace | string | `""` | Gateway namespace (optional) |
| gateway.hostnames | list | `["chart-example.local"]` | Hostnames to match |
| gateway.rules | list | `[]` | Path rules (default: catch-all) |
| livenessProbe.httpGet.path | string | `"/health"` |  |
| livenessProbe.httpGet.port | string | `"http"` |  |
| log.level | string | `"Info"` |  |
| nameOverride | string | `""` |  |
| podAnnotations | object | `{}` |  |
| podLabels | object | `{}` |  |
| podSecurityContext.fsGroup | int | `65534` |  |
| postgresql.auth.database | string | `"kube-app"` |  |
| postgresql.auth.enablePostgresUser | bool | `false` |  |
| postgresql.auth.host | string | `"postgresql"` |  |
| postgresql.auth.password | string | `""` |  |
| postgresql.auth.port | int | `5432` |  |
| postgresql.auth.username | string | `""` |  |
| postgresql.backup.cronjob.containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| postgresql.backup.cronjob.containerSecurityContext.capabilities.drop[0] | string | `"ALL"` |  |
| postgresql.backup.cronjob.containerSecurityContext.enabled | bool | `false` |  |
| postgresql.backup.cronjob.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| postgresql.backup.cronjob.containerSecurityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| postgresql.backup.cronjob.podSecurityContext.enabled | bool | `false` |  |
| postgresql.enabled | bool | `true` |  |
| postgresql.metrics.containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| postgresql.metrics.containerSecurityContext.capabilities.drop[0] | string | `"ALL"` |  |
| postgresql.metrics.containerSecurityContext.enabled | bool | `false` |  |
| postgresql.metrics.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| postgresql.metrics.containerSecurityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| postgresql.primary.containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| postgresql.primary.containerSecurityContext.capabilities.drop[0] | string | `"ALL"` |  |
| postgresql.primary.containerSecurityContext.enabled | bool | `false` |  |
| postgresql.primary.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| postgresql.primary.containerSecurityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| postgresql.primary.podSecurityContext.enabled | bool | `false` |  |
| postgresql.readReplicas.containerSecurityContext.allowPrivilegeEscalation | bool | `false` |  |
| postgresql.readReplicas.containerSecurityContext.capabilities.drop[0] | string | `"ALL"` |  |
| postgresql.readReplicas.containerSecurityContext.enabled | bool | `false` |  |
| postgresql.readReplicas.containerSecurityContext.runAsNonRoot | bool | `true` |  |
| postgresql.readReplicas.containerSecurityContext.seccompProfile.type | string | `"RuntimeDefault"` |  |
| postgresql.readReplicas.podSecurityContext.enabled | bool | `false` |  |
| postgresql.tls.enabled | bool | `false` |  |
| readinessProbe.httpGet.path | string | `"/health"` |  |
| readinessProbe.httpGet.port | string | `"http"` |  |
| replicaCount | int | `1` |  |
| resources.limits.cpu | string | `"100m"` |  |
| resources.limits.memory | string | `"128Mi"` |  |
| resources.requests.cpu | string | `"100m"` |  |
| resources.requests.memory | string | `"128Mi"` |  |
| securityContext.capabilities.drop[0] | string | `"ALL"` |  |
| securityContext.readOnlyRootFilesystem | bool | `true` |  |
| securityContext.runAsNonRoot | bool | `true` |  |
| securityContext.runAsUser | int | `65534` |  |
| service.port | int | `80` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.automount | bool | `true` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `""` |  |
| volumeMounts | list | `[]` |  |
| volumes | list | `[]` |  |
