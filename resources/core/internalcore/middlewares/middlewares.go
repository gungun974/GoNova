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
