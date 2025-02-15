package make_entity_template

import (
	_ "embed"
	"text/template"
)

//go:embed entity.go.tmpl
var entityGoFile []byte

var EntityGoTemplate = template.Must(template.New("make:entity.go").Parse(string(entityGoFile)))
