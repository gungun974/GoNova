package make_presenter_template

import (
	_ "embed"
	"text/template"
)

//go:embed presenter.go.tmpl
var presenterGoFile []byte

var PresenterGoTemplate = template.Must(
	template.New("make:presenter.go").Parse(string(presenterGoFile)),
)
