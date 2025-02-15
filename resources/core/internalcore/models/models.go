package models_template

import (
	_ "embed"
	"text/template"
)

//go:embed api.go.tmpl
var apiGoFile []byte

var ApiGoTemplate = template.Must(template.New("api.go").Parse(string(apiGoFile)))

//go:embed file.go.tmpl
var fileGoFile []byte

var FileGoTemplate = template.Must(template.New("file.go").Parse(string(fileGoFile)))

//go:embed html.go.tmpl
var htmlGoFile []byte

var HtmlGoTemplate = template.Must(template.New("html.go").Parse(string(htmlGoFile)))

//go:embed image.go.tmpl
var imageGoFile []byte

var ImageGoTemplate = template.Must(template.New("image.go").Parse(string(imageGoFile)))

//go:embed json.go.tmpl
var jsonGoFile []byte

var JsonGoTemplate = template.Must(template.New("json.go").Parse(string(jsonGoFile)))

//go:embed pdf.go.tmpl
var pdfGoFile []byte

var PdfGoTemplate = template.Must(template.New("pdf.go").Parse(string(pdfGoFile)))

//go:embed plain.go.tmpl
var plainGoFile []byte

var PlainGoTemplate = template.Must(template.New("plain.go").Parse(string(plainGoFile)))

//go:embed reader.go.tmpl
var readerGoFile []byte

var ReaderGoTemplate = template.Must(template.New("reader.go").Parse(string(readerGoFile)))

//go:embed redirect.go.tmpl
var redirectGoFile []byte

var RedirectGoTemplate = template.Must(template.New("redirect.go").Parse(string(redirectGoFile)))
