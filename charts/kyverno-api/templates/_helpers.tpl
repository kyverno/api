{{/* vim: set filetype=mustache: */}}

{{- define "kyverno-api.chartVersion" -}}
{{- .Chart.Version | replace "+" "_" -}}
{{- end -}}

{{- define "kyverno-api.labels" -}}
{{- $customLabels := .Values.labels | default dict -}}
{{- if not (hasKey $customLabels "helm.sh/chart") }}
helm.sh/chart: {{ .Chart.Name }}-{{ include "kyverno-api.chartVersion" . }}
{{- end }}
{{- if not (hasKey $customLabels "app.kubernetes.io/component") }}
app.kubernetes.io/component: kyverno-api
{{- end }}
{{- if not (hasKey $customLabels "app.kubernetes.io/instance") }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
{{- if not (hasKey $customLabels "app.kubernetes.io/managed-by") }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}
{{- if not (hasKey $customLabels "app.kubernetes.io/part-of") }}
app.kubernetes.io/part-of: kyverno-api
{{- end }}
{{- if not (hasKey $customLabels "app.kubernetes.io/version") }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
{{- with .Values.labels }}
{{- tpl (toYaml .) $ -}}
{{- end }}
{{- end -}}

{{- define "kyverno-api.annotations" -}}
{{- tpl (toYaml .Values.annotations) . -}}
{{- end -}}
