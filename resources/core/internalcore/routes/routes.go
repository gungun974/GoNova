package routes_template

import (
	_ "embed"
	"text/template"
)

//go:embed file.go.tmpl
var fileGoFile []byte

var FileGoTemplate = template.Must(template.New("file.go").Parse(string(fileGoFile)))

//go:embed handle_http_error.go.tmpl
var handleHttpErrorGoFile []byte

var HandleHttpErrorGoTemplate = template.Must(template.New("handle_http_error.go").Parse(string(handleHttpErrorGoFile)))

//go:embed home.go.tmpl
var homeGoFile []byte

var HomeGoTemplate = template.Must(template.New("home.go").Parse(string(homeGoFile)))

//go:embed routes.go.tmpl
var routesGoFile []byte

var RoutesGoTemplate = template.Must(template.New("routes.go").Parse(string(routesGoFile)))
