package injector

import (
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/gungun974/gonova/internal/logger"
)

func InjectContainerDatabase(path string) {
	f, err := decorator.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	addImport(f, "github.com/jmoiron/sqlx", "")

	found := false

	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*dst.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Name.Name != "NewContainer" {
			continue
		}

		found = true

		funcDecl.Type.Params.List = append(
			[]*dst.Field{
				{
					Names: []*dst.Ident{{Name: "db"}},
					Type: &dst.StarExpr{
						X: &dst.SelectorExpr{
							Sel: dst.NewIdent("DB"),
							X:   dst.NewIdent("sqlx"),
						},
					},
				},
			},
			funcDecl.Type.Params.List...,
		)
	}

	if !found {
		logger.InjectorLogger.Fatal("Failed to find `NewContainer` in container.go")
	}

	err = writeDSTFileToPath(f, path)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}
}
