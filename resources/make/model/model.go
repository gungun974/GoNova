package make_model_template

import (
	_ "embed"
	"text/template"
)

//go:embed model.go.tmpl
var modelGoFile []byte

var ModelGoTemplate = template.Must(template.New("make:model.go").Parse(string(modelGoFile)))
