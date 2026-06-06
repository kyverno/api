{{/* vim: set filetype=mustache: */}}

{{- define "kyverno-api.chartVersion" -}}
{{- .Chart.Version | replace "+" "_" -}}
{{- end -}}

{{- define "kyverno-api.labels" -}}
helm.sh/chart: crds-{{ include "kyverno-api.chartVersion" . }}
app.kubernetes.io/component: kyverno-api
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: kyverno-api
app.kubernetes.io/version: {{ include "kyverno-api.chartVersion" . }}
{{- with .Values.labels }}
{{ tpl (toYaml .) $ }}
{{- end }}
{{- end -}}

{{- define "kyverno-api.annotations" -}}
{{- tpl (toYaml .Values.annotations) . -}}
{{- end -}}
