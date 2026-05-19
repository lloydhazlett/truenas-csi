{{- define "truenas-csi.namespace" -}}
{{- .Values.namespace }}
{{- end }}

{{- define "truenas-csi.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}
