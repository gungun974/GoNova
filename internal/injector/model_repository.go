package injector

import (
	"bytes"
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	make_repository_template "github.com/gungun974/gonova/resources/make/repository"
	"github.com/jinzhu/inflection"
)

func InjectModelInRepository(path string, repositoryName string, model analyzer.AnalyzedModel) {
	projectName, err := utils.GetGoModName(".")
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	projectGlobalTemplateConfig := struct {
		ProjectName string

		ModelName  string
		ModelsName string

		EntityName   string
		EntitiesName string

		RepositoryName string

		IdType string
	}{
		ProjectName: projectName,

		ModelName:  helpers.CapitalizeFirstLetter(model.Name),
		ModelsName: inflection.Plural(helpers.CapitalizeFirstLetter(model.Name)),

		EntityName:   helpers.CapitalizeFirstLetter(model.Entity.Name),
		EntitiesName: inflection.Plural(helpers.CapitalizeFirstLetter(model.Entity.Name)),

		RepositoryName: repositoryName,

		IdType: "int",
	}

	var buffer bytes.Buffer
	err = make_repository_template.RepositoryModelGoTemplate.
		Execute(&buffer, projectGlobalTemplateConfig)
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	newDst, err := decorator.Parse(buffer.String())
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	f, err := decorator.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	f.Decls = append(f.Decls, newDst.Decls...)

	newDecls := make([]dst.Decl, 0, len(f.Decls)+1)

	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*dst.GenDecl); ok && genDecl.Tok == token.IMPORT {
			newDecls = append(newDecls, decl)
		}
	}

	newDecls = append(newDecls, &dst.GenDecl{
		Tok: token.VAR,
		Specs: []dst.Spec{
			&dst.ValueSpec{
				Names: []*dst.Ident{
					dst.NewIdent(model.Entity.Name + "NotFoundError"),
				},
				Values: []dst.Expr{
					&dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("errors"),
							Sel: dst.NewIdent("New"),
						},
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: "\"" + model.Entity.Name + " is not found\"",
							},
						},
					},
				},
			},
		},
	})

	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*dst.GenDecl); !ok || genDecl.Tok != token.IMPORT {
			newDecls = append(newDecls, decl)
		}
	}

	f.Decls = newDecls

	err = writeDSTFileToPath(f, path)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}
}
