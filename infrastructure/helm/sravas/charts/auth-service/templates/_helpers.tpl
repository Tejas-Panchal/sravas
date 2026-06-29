{{- define "sravas.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "sravas.fullname" -}}
{{- if .Values.global.namePrefix }}{{ .Values.global.namePrefix }}{{- end }}{{ include "sravas.name" . }}
{{- end }}

{{- define "sravas.labels" -}}
helm.sh/chart: {{ include "sravas.name" . }}
app.kubernetes.io/name: {{ include "sravas.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "sravas.selectorLabels" -}}
app.kubernetes.io/name: {{ include "sravas.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
