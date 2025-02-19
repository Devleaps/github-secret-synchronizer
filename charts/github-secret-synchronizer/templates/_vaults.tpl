{{/*
JSON env block
*/}}
{{- define "github-secret-synchronizer.json-env" -}}
- name: VAULT_TYPE
  value: "json"
- name: JSON_VAULT_FILE_PATH
  value: secrets.json
{{- end }}

{{/*
JSON volumemount block
*/}}
{{- define "github-secret-synchronizer.json-volumeMount" -}}
- name: json-secrets
  mountPath: /secrets.json
  subPath: secrets.json
{{- end }}

{{/*
JSON volume block
*/}}
{{- define "github-secret-synchronizer.json-volume" -}}
- name: json-secrets
  secret:
    secretName: {{ include "github-secret-synchronizer.fullname" . }}-vault
{{- end }}

{{/*
YAML env block
*/}}
{{- define "github-secret-synchronizer.yaml-env" -}}
- name: VAULT_TYPE
  value: "yaml"
- name: YAML_VAULT_FILE_PATH
  value: secrets.yaml
{{- end }}

{{/*
YAML volumemount block
*/}}
{{- define "github-secret-synchronizer.yaml-volumeMount" -}}
- name: yaml-secrets
  mountPath: /secrets.yaml
  subPath: secrets.yaml
{{- end }}

{{/*
YAML volume block
*/}}
{{- define "github-secret-synchronizer.yaml-volume" -}}
- name: yaml-secrets
  secret:
    secretName: {{ include "github-secret-synchronizer.fullname" . }}-vault
{{- end }}

{{/*
Azure env block
*/}}
{{- define "github-secret-synchronizer.azure-env" -}}
- name: VAULT_TYPE
  value: "azure"
- name: AZURE_KEYVAULT_URL
  valueFrom:
    secretKeyRef:
      name: {{ default (printf "%s-%s" (include "github-secret-synchronizer.fullname" .) "vault") .Values.synchronizer.azure.existingSecretName }}
      key: keyvaultURL
- name: AZURE_CLIENT_ID
  valueFrom:
    secretKeyRef:
      name: {{ default (printf "%s-%s" (include "github-secret-synchronizer.fullname" .) "vault") .Values.synchronizer.azure.existingSecretName }}
      key: clientID
- name: AZURE_TENANT_ID
  valueFrom:
    secretKeyRef:
      name: {{ default (printf "%s-%s" (include "github-secret-synchronizer.fullname" .) "vault") .Values.synchronizer.azure.existingSecretName }}
      key: tenantID
- name: AZURE_CLIENT_SECRET
  valueFrom:
    secretKeyRef:
      name: {{ default (printf "%s-%s" (include "github-secret-synchronizer.fullname" .) "vault") .Values.synchronizer.azure.existingSecretName }}
      key: clientSecret
{{- end }}

{{/*
AWS env block
*/}}
{{- define "github-secret-synchronizer.aws-env" -}}
- name: VAULT_TYPE
  value: "aws"
- name: AWS_ACCESS_KEY_ID
  valueFrom:
    secretKeyRef:
      name: {{ default (printf "%s-%s" (include "github-secret-synchronizer.fullname" .) "vault") .Values.synchronizer.aws.existingSecretName }}
      key: accessKeyID
- name: AWS_SECRET_ACCESS_KEY
  valueFrom:
    secretKeyRef:
      name: {{ default (printf "%s-%s" (include "github-secret-synchronizer.fullname" .) "vault") .Values.synchronizer.aws.existingSecretName }}
      key: secretAccessKey
- name: AWS_REGION
  valueFrom:
    secretKeyRef:
      name: {{ default (printf "%s-%s" (include "github-secret-synchronizer.fullname" .) "vault") .Values.synchronizer.aws.existingSecretName }}
      key: region
{{- end }}
