package home_template

import (
	_ "embed"
	"text/template"
)

//go:embed index.templ.tmpl
var indexTemplFile []byte

var IndexTemplTemplate = template.Must(template.New("index.templ").Parse(string(indexTemplFile)))
