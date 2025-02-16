package injector

import (
	"bytes"
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
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
	err = make_model_template.InjectModelGoTemplate.Execute(&buffer, projectGlobalTemplateConfig)
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	newDst, err := decorator.Parse(buffer.String())
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	dstutil.Apply(newDst, nil, func(c *dstutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *dst.CompositeLit:
			selectorExpr, ok := x.Type.(*dst.SelectorExpr)
			if !ok {
				return true
			}

			mod, ok := selectorExpr.X.(*dst.Ident)
			if !ok {
				return true
			}

			if mod.Name != "entities" || selectorExpr.Sel.Name != projectGlobalTemplateConfig.EntityName {
				return true
			}

			for _, field := range entity.Fields {
				if !field.Exported() {
					continue
				}
				x.Elts = append(x.Elts, &dst.KeyValueExpr{
					Key:   dst.NewIdent(field.Name()),
					Value: dst.NewIdent(helpers.LowerFirstLetter(field.Name())),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{
							After: dst.NewLine,
						},
					},
				})
			}
		}

		return true
	})

	f, err := decorator.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

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
