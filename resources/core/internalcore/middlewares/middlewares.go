package middlewares_template

import (
	_ "embed"
	"text/template"
)

//go:embed conditional_logger.go.tmpl
var conditionalLoggerGoFile []byte

var ConditionalLoggerGoTemplate = template.Must(template.New("conditional_logger.go").Parse(string(conditionalLoggerGoFile)))

//go:embed request_info_middleware.go.tmpl
var requestInfoMiddlewareGoFile []byte

var RequestInfoMiddlewareGoTemplate = template.Must(template.New("request_info_middleware.go").Parse(string(requestInfoMiddlewareGoFile)))

//go:embed vite.go.tmpl
var viteGoFile []byte

var ViteGoTemplate = template.Must(template.New("vite.go").Parse(string(viteGoFile)))

//go:embed vite_debug.go.tmpl
var viteDebugGoFile []byte

var ViteDebugGoTemplate = template.Must(template.New("vite_debug.go").Parse(string(viteDebugGoFile)))
