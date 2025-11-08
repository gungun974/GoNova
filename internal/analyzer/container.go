package analyzer

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/gungun974/gonova/internal/logger"
)

type AnalyzedContainer struct {
	Dependencies []AnalyzedDependency
}

func AnalyzeProjectContainer(controllers []AnalyzedController) AnalyzedContainer {
	projectPath := "."

	containerFilePath := filepath.Join(projectPath, "/internal/container.go")

	f, err := decorator.ParseFile(token.NewFileSet(), containerFilePath, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	dependencies := []AnalyzedDependency{}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*dst.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*dst.TypeSpec)
			if !ok {
				continue
			}

			if typeSpec.Name == nil {
				continue
			}

			if typeSpec.Name.Name != "Container" {
				continue
			}

			structType, ok := typeSpec.Type.(*dst.StructType)
			if !ok {
				continue
			}

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				expr := field.Type

				starExpr, ok := field.Type.(*dst.StarExpr)
				if ok {
					expr = starExpr.X
				}

				selectorExpr, ok := expr.(*dst.SelectorExpr)
				if !ok {
					continue
				}

				identX, ok := selectorExpr.X.(*dst.Ident)
				if !ok {
					continue
				}

				if strings.Contains(identX.Name, "controllers") {
					for _, controller := range controllers {
						if controller.Name == selectorExpr.Sel.Name {
							dependencies = append(dependencies, &controller)
							break
						}
					}
				}
			}
		}
	}

	return AnalyzedContainer{
		Dependencies: dependencies,
	}
}
