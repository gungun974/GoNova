package make_usecase_template

import (
	_ "embed"
	"text/template"
)

//go:embed usecase.go.tmpl
var usecaseGoFile []byte

var UsecaseGoTemplate = template.Must(template.New("make:usecase.go").Parse(string(usecaseGoFile)))

//go:embed usecase_get_example.go.tmpl
var usecaseGetExampleGoFile []byte

var UsecaseGetExampleGoTemplate = template.Must(
	template.New("make:usecase_get_example.go").Parse(string(usecaseGetExampleGoFile)),
)
