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
	make_model_template "github.com/gungun974/gonova/resources/make/model"
	"github.com/jinzhu/inflection"
)

func InjectModelNewModel(path string, entity analyzer.AnalyzedEntity) {
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
	}{
		ProjectName: projectName,

		ModelName:  helpers.CapitalizeFirstLetter(entity.Name) + "Model",
		ModelsName: inflection.Plural(helpers.CapitalizeFirstLetter(entity.Name) + "Model"),

		EntityName:   helpers.CapitalizeFirstLetter(entity.Name),
		EntitiesName: inflection.Plural(helpers.CapitalizeFirstLetter(entity.Name)),
	}

	var buffer bytes.Buffer
	make_model_template.InjectModelGoTemplate.Execute(&buffer, projectGlobalTemplateConfig)

	newDst, err := decorator.Parse(buffer.String())
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	f, err := decorator.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	// addImport(f, projectName+"/internal/layers/domain/entities", "")

	f.Decls = append(f.Decls, newDst.Decls...)

	newDecls := make([]dst.Decl, 0, len(f.Decls))

	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*dst.GenDecl); ok && genDecl.Tok == token.IMPORT {
			newDecls = append(newDecls, decl)
		}
	}

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
