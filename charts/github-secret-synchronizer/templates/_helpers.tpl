{{/*
Expand the name of the chart.
*/}}
{{- define "github-secrets-synchronizer.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "github-secrets-synchronizer.fullname" -}}
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
{{- define "github-secrets-synchronizer.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "github-secrets-synchronizer.labels" -}}
helm.sh/chart: {{ include "github-secrets-synchronizer.chart" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
GitHub env block
*/}}
{{- define "github-secrets-synchronizer.github-env" -}}
- name: GITHUB_APP_ID
  valueFrom:
    secretKeyRef:
        name: {{ default (printf "%s-%s" (include "github-secrets-synchronizer.fullname" .) "github") .Values.github.existingSecretName }}
        key: appID
- name: GITHUB_APP_INSTALLATION_ID
  valueFrom:
    secretKeyRef:
        name: {{ default (printf "%s-%s" (include "github-secrets-synchronizer.fullname" .) "github") .Values.github.existingSecretName }}
        key: appInstallationID
- name: GITHUB_APP_PRIVATE_KEY
  valueFrom:
    secretKeyRef:
        name: {{ default (printf "%s-%s" (include "github-secrets-synchronizer.fullname" .) "github") .Values.github.existingSecretName }}
        key: appPrivateKey
- name: GITHUB_ORG_NAME
  valueFrom:
    secretKeyRef:
        name: {{ default (printf "%s-%s" (include "github-secrets-synchronizer.fullname" .) "github") .Values.github.existingSecretName }}
        key: organization
{{- end }}

{{/*
Defaults for vault
*/}}
{{- define "github-secrets-synchronizer.defaults" -}}
{{- if .Values.defaults.type }}
- name: VAULT_DEFAULT_TYPE
  value: {{ .Values.defaults.type }}
{{- end }}
{{- if .Values.defaults.visibility }}
- name: VAULT_DEFAULT_VISIBILITY
  value: {{ .Values.defaults.visibility }}
{{- end }}
{{- if .Values.defaults.repositories }}
- name: VAULT_DEFAULT_REPOSITORIES
  value: {{ .Values.defaults.repositories }}
{{- end }}
{{- end }}
