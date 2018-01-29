{{ range . -}}
{{escape .first_name}} {{escape .last_name}}
{{- end -}}
