package injector

import (
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
)

func injectUsecaseDependency(path string, dependencyPath string, dependencyType string, usecase analyzer.AnalyzedUsecase, dependencyName string) {
	projectName, err := utils.GetGoModName(".")
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	f, err := decorator.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	addImport(f, projectName+"/"+dependencyPath+"/"+dependencyType, "")

	foundStruct := false
	foundFunction := false

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*dst.GenDecl)
		if ok {
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*dst.TypeSpec)
				if !ok {
					continue
				}

				if typeSpec.Name == nil {
					continue
				}

				if typeSpec.Name.Name != usecase.Name {
					continue
				}

				structType, ok := typeSpec.Type.(*dst.StructType)
				if !ok {
					continue
				}

				foundStruct = true
				skip := false

				for _, field := range structType.Fields.List {
					if len(field.Names) != 0 && field.Names[0].Name == helpers.LowerFirstLetter(dependencyName) {
						skip = true
						break
					}
				}

				if skip {
					continue
				}

				structType.Fields.List = append(structType.Fields.List, &dst.Field{
					Names: []*dst.Ident{
						dst.NewIdent(helpers.LowerFirstLetter(dependencyName)),
					},
					Type: &dst.SelectorExpr{
						X:   dst.NewIdent(dependencyType),
						Sel: dst.NewIdent(dependencyName),
					},
					Decs: dst.FieldDecorations{
						NodeDecs: dst.NodeDecs{
							After: dst.NewLine,
						},
					},
				})

			}

			continue
		}

		funcDecl, ok := decl.(*dst.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Name.Name != "New"+helpers.CapitalizeFirstLetter(usecase.Name) {
			continue
		}

		foundFunction = true

		skip := false

		for _, field := range funcDecl.Type.Params.List {
			for _, ident := range field.Names {
				if ident.Name != helpers.LowerFirstLetter(dependencyName) {
					continue
				}

				skip = true

				break
			}
		}

		if !skip {
			funcDecl.Type.Params.List = append(
				funcDecl.Type.Params.List,
				&dst.Field{
					Names: []*dst.Ident{
						dst.NewIdent(helpers.LowerFirstLetter(dependencyName)),
					},
					Type: &dst.SelectorExpr{
						X:   dst.NewIdent(dependencyType),
						Sel: dst.NewIdent(dependencyName),
					},
					Decs: dst.FieldDecorations{
						NodeDecs: dst.NodeDecs{
							Before: dst.NewLine,
							After:  dst.NewLine,
						},
					},
				},
			)
		}

		var returnCompositeLit *dst.CompositeLit

		for _, stmt := range funcDecl.Body.List {
			returnStmt, ok := stmt.(*dst.ReturnStmt)
			if ok {
				for _, expr := range returnStmt.Results {
					compositeLit, ok := expr.(*dst.CompositeLit)
					if !ok {
						continue
					}
					ident, ok := compositeLit.Type.(*dst.Ident)
					if !ok {
						continue
					}

					if ident.Name != usecase.Name {
						continue
					}

					returnCompositeLit = compositeLit
					break
				}
				continue
			}
		}

		if returnCompositeLit == nil {
			logger.InjectorLogger.Fatalf("Failed to find function `%s` with valid return in %s", "New"+helpers.CapitalizeFirstLetter(usecase.Name), path)
		}

		skip = false

		for _, elt := range returnCompositeLit.Elts {
			ident, ok := elt.(*dst.Ident)
			if !ok {
				continue
			}

			if ident.Name != helpers.LowerFirstLetter(usecase.Name) {
				continue
			}

			skip = true

			break
		}

		if !skip {
			returnCompositeLit.Elts = append(
				returnCompositeLit.Elts,
				&dst.Ident{
					Name: helpers.LowerFirstLetter(dependencyName),
					Decs: dst.IdentDecorations{
						NodeDecs: dst.NodeDecs{
							Before: dst.NewLine,
							After:  dst.NewLine,
						},
					},
				},
			)
		}

	}

	if !foundStruct {
		logger.InjectorLogger.Fatalf("Failed to find struct `%s` in %s", usecase.Name, path)
	}

	if !foundFunction {
		logger.InjectorLogger.Fatalf("Failed to find function `%s` in %s", "New"+helpers.CapitalizeFirstLetter(usecase.Name), path)
	}

	err = writeDSTFileToPath(f, path)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}
}

func InjectUsecaseRepository(path string, usecase analyzer.AnalyzedUsecase, repository analyzer.AnalyzedRepository) {
	injectUsecaseDependency(path, "internal/layers/data", "repositories", usecase, repository.Name)
}

func InjectUsecasePresenter(path string, usecase analyzer.AnalyzedUsecase, presenter analyzer.AnalyzedPresenter) {
	injectUsecaseDependency(path, "internal/layers/presentation", "presenters", usecase, presenter.Name)
}
