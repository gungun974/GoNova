package errors_page_template

import (
	_ "embed"
	"text/template"
)

//go:embed errors.templ.tmpl
var errorsTemplFile []byte

var ErrorsTemplTemplate = template.Must(template.New("errors.templ").Parse(string(errorsTemplFile)))
