package server_template

import (
	_ "embed"
	"text/template"
)

//go:embed main.go.tmpl
var mainGoFile []byte

var MainGoTemplate = template.Must(template.New("main.go").Parse(string(mainGoFile)))
