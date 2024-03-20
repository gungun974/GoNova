package scripts_template

import (
	_ "embed"
	"text/template"
)

//go:embed core.ts.tmpl
var coreTSFile []byte

var CoreTSTemplate = template.Must(template.New("core.ts").Parse(string(coreTSFile)))

//go:embed htmx/core.ts.tmpl
var htmxCoreTSFile []byte

var HtmxCoreTSTemplate = template.Must(template.New("htmx/core.ts").Parse(string(htmxCoreTSFile)))

//go:embed htmx/extensions.ts.tmpl
var htmxExtensionsTSFile []byte

var HtmxExtensionsTSTemplate = template.Must(template.New("htmx/extensions.ts").Parse(string(htmxExtensionsTSFile)))

//go:embed htmx/index.ts.tmpl
var htmxIndexTSFile []byte

var HtmxIndexTSTemplate = template.Must(template.New("htmx/index.ts").Parse(string(htmxIndexTSFile)))
