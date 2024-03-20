package css_template

import (
	_ "embed"
	"text/template"
)

//go:embed main.css.tmpl
var mainCssFile []byte

var MainCssTemplate = template.Must(template.New("main.css").Parse(string(mainCssFile)))

//go:embed tailwind.css.tmpl
var tailwindCssFile []byte

var TailwindCssTemplate = template.Must(template.New("tailwind.css").Parse(string(tailwindCssFile)))
