{{ range . -}}
{{ pluralize_with_size .fruit (int .count )}}
{{ end -}}
