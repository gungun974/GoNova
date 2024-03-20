package layouts_template

import (
	_ "embed"
	"text/template"
)

//go:embed app_layout.templ.tmpl
var appLayoutTemplFile []byte

var AppLayoutTemplTemplate = template.Must(template.New("app_layout.templ").Parse(string(appLayoutTemplFile)))

//go:embed core.templ.tmpl
var coreTemplFile []byte

var CoreTemplTemplate = template.Must(template.New("core.templ").Parse(string(coreTemplFile)))
