package entities_template

import (
	_ "embed"
	"text/template"
)

//go:embed error.go.tmpl
var errorGoFile []byte

var ErrorGoTemplate = template.Must(template.New("error.go").Parse(string(errorGoFile)))
