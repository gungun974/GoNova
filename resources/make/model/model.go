package make_model_template

import (
	_ "embed"
	"text/template"
)

//go:embed blank_model.go.tmpl
var blankModelGoFile []byte

var BlankModelGoTemplate = template.Must(
	template.New("make:blank_model.go").Parse(string(blankModelGoFile)),
)

//go:embed inject_model.go.tmpl
var injectjodelGoFile []byte

var InjectModelGoTemplate = template.Must(
	template.New("make:inject_model.go").Parse(string(injectjodelGoFile)),
)
