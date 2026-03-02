{{/*
Expand the name of the chart.
*/}}
{{- define "kube-app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kube-app.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "kube-app.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kube-app.labels" -}}
helm.sh/chart: {{ include "kube-app.chart" . }}
{{ include "kube-app.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kube-app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kube-app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "kube-app.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "kube-app.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Database credentials (used only when cnpg.enabled is false, for dbsecret)
*/}}
{{- define "dbUser" -}}
{{- default "kube-app" (index .Values "cnpg" "auth" "username") -}}
{{- end -}}
{{- define "dbPassword" -}}
{{- default "kube-app" (index .Values "cnpg" "auth" "password") -}}
{{- end -}}
{{- define "staticToken" -}}
{{- ternary "kube-app" .Values.api.token (empty .Values.api.token) -}}
{{- end -}}
