package make_repository_template

import (
	_ "embed"
	"text/template"

	"github.com/gungun974/gonova/internal/helpers"
)

//go:embed empty_repository.go.tmpl
var emptyRepositoryGoFile []byte

var EmptyRepositoryGoTemplate = template.Must(
	template.New("make:empty_repository.go").Parse(string(emptyRepositoryGoFile)),
)

//go:embed repository_model.go.tmpl
var repositoryModelGoFile []byte

var RepositoryModelGoTemplate = template.Must(
	template.New("make:repository_model.go").Funcs(
		template.FuncMap{
			"lowerFirst": helpers.LowerFirstLetter,
		},
	).Parse(string(repositoryModelGoFile)),
)
