package make_repository_template

import (
	_ "embed"
	"text/template"
)

//go:embed empty_repository.go.tmpl
var emptyRepositoryGoFile []byte

var EmptyRepositoryGoTemplate = template.Must(
	template.New("make:empty_repository.go").Parse(string(emptyRepositoryGoFile)),
)
