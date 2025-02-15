package make_router_template

import (
	_ "embed"
	"text/template"
)

//go:embed route.go.tmpl
var routeGoFile []byte

var RouteGoTemplate = template.Must(template.New("make:route.go").Parse(string(routeGoFile)))
