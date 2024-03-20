package context_template

import (
	_ "embed"
	"text/template"
)

//go:embed key.go.tmpl
var keyGoFile []byte

var KeyGoTemplate = template.Must(template.New("key.go").Parse(string(keyGoFile)))

//go:embed request_info.go.tmpl
var requestInfoGoFile []byte

var RequestInfoGoTemplate = template.Must(template.New("request_info.go").Parse(string(requestInfoGoFile)))
