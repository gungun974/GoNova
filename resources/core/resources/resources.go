package resources_template

import (
	_ "embed"
	"text/template"
)

//go:embed dom.d.ts.tmpl
var domDTSFile []byte

var DomDTSTemplate = template.Must(template.New("dom.d.ts").Parse(string(domDTSFile)))

//go:embed main.ts.tmpl
var mainTSFile []byte

var MainTSTemplate = template.Must(template.New("main.ts").Parse(string(mainTSFile)))
