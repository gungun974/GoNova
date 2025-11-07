package make_controller_template

import (
	_ "embed"
	"text/template"
)

//go:embed controller.go.tmpl
var controllerGoFile []byte

var ControllerGoTemplate = template.Must(
	template.New("make:controller.go").Parse(string(controllerGoFile)),
)
