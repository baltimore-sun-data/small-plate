{{ range groupby "last_name" . -}}
{{.Key}}s:
{{- range .Items }}
+ {{.first_name}}
{{- end }}
{{ end -}}
