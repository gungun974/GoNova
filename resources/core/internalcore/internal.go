package internal_template

import (
	_ "embed"
	"text/template"
)

//go:embed container.go.tmpl
var containerGoFile []byte

var ContainerGoTemplate = template.Must(template.New("container.go").Parse(string(containerGoFile)))
