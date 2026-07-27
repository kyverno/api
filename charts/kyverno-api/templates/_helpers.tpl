{{/* vim: set filetype=mustache: */}}

{{- define "kyverno-api.chartVersion" -}}
{{- .Chart.Version | replace "+" "_" -}}
{{- end -}}

{{- define "kyverno-api.labels" -}}
{{- with .Values.labels }}
{{- tpl (toYaml .) $ }}
{{- end }}
{{- end -}}

{{- define "kyverno-api.annotations" -}}
{{- with .Values.annotations }}
{{- tpl (toYaml .) $ }}
{{- end }}
{{- end -}}
