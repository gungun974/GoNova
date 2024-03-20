package logger_template

import (
	_ "embed"
	"text/template"
)

//go:embed logger.go.tmpl
var loggerGoFile []byte

var LoggerGoTemplate = template.Must(template.New("logger.go").Parse(string(loggerGoFile)))
