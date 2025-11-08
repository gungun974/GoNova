package injector

import (
	"go/parser"
	"go/token"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
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

func InjectContainerController(path string, controller analyzer.AnalyzedController) {
	projectName, err := utils.GetGoModName(".")
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	f, err := decorator.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	addImport(f, projectName+"/internal/layers/presentation/controllers", "")

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

				if typeSpec.Name.Name != "Container" {
					continue
				}

				structType, ok := typeSpec.Type.(*dst.StructType)
				if !ok {
					continue
				}

				foundStruct = true
				skip := false

				for _, field := range structType.Fields.List {
					if len(field.Names) != 0 && field.Names[0].Name == helpers.CapitalizeFirstLetter(controller.Name) {
						skip = true
						break
					}
				}

				if skip {
					continue
				}

				structType.Fields.List = append(structType.Fields.List, &dst.Field{
					Names: []*dst.Ident{
						dst.NewIdent(helpers.CapitalizeFirstLetter(controller.Name)),
					},
					Type: &dst.SelectorExpr{
						X:   dst.NewIdent("controllers"),
						Sel: dst.NewIdent(controller.Name),
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

		if funcDecl.Name.Name != "NewContainer" {
			continue
		}

		foundFunction = true

		var containerName string
		var insertPos int

		for i, stmt := range funcDecl.Body.List {
			returnStmt, ok := stmt.(*dst.ReturnStmt)
			if ok {
				for _, expr := range returnStmt.Results {
					ident, ok := expr.(*dst.Ident)
					if ok {
						containerName = ident.Name
						insertPos = i
						break
					}
				}
				continue
			}
		}

		if containerName == "" {
			logger.InjectorLogger.Fatal("Failed to find Container return variable in container.go")
		}

		skip := false

		for _, stmt := range funcDecl.Body.List {
			assignStmt, ok := stmt.(*dst.AssignStmt)
			if !ok {
				continue
			}
			for _, expr := range assignStmt.Lhs {
				selectorExpr, ok := expr.(*dst.SelectorExpr)
				if !ok {
					continue
				}

				identX, ok := selectorExpr.X.(*dst.Ident)
				if !ok {
					continue
				}

				if identX.Name != "container" {
					continue
				}

				if selectorExpr.Sel.Name != helpers.CapitalizeFirstLetter(controller.Name) {
					continue
				}

				skip = true

				break
			}
		}

		if skip {
			continue
		}

		funcDecl.Body.List = append(funcDecl.Body.List[:insertPos], append([]dst.Stmt{
			&dst.AssignStmt{
				Lhs: []dst.Expr{
					&dst.SelectorExpr{
						X:   dst.NewIdent("container"),
						Sel: dst.NewIdent(helpers.CapitalizeFirstLetter(controller.Name)),
					},
				},
				Rhs: []dst.Expr{
					&dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("controllers"),
							Sel: dst.NewIdent("New" + helpers.CapitalizeFirstLetter(controller.Name)),
						},
					},
				},
				Tok: token.ASSIGN,
			},
		}, funcDecl.Body.List[insertPos:]...)...)
	}

	if !foundStruct {
		logger.InjectorLogger.Fatal("Failed to find struct `Container` in container.go")
	}

	if !foundFunction {
		logger.InjectorLogger.Fatal("Failed to find function `NewContainer` in container.go")
	}

	err = writeDSTFileToPath(f, path)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}
}

func InjectContainerDependencies(path string, target analyzer.AnalyzedDependency) {
	f, err := decorator.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	dependencies := target.GetDependencies()

	for _, dependency := range dependencies {
		if dependency == nil {
			continue
		}
		if strings.Contains(dependency.GetName(), "Usecase") {
			addImport(f, dependency.GetPkgPath(), dependency.GetImportName())
		} else {
			addImport(f, dependency.GetPkgPath(), "")
		}
	}

	foundFunction := false

	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*dst.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Name.Name != "NewContainer" {
			continue
		}

		foundFunction = true
		var firstNewMethod *int

		for i, stmt := range funcDecl.Body.List {
			assignStmt, ok := stmt.(*dst.AssignStmt)
			if !ok {
				continue
			}
			for _, expr := range assignStmt.Rhs {
				callExpr, ok := expr.(*dst.CallExpr)
				if !ok {
					continue
				}

				selectorExpr, ok := callExpr.Fun.(*dst.SelectorExpr)
				if !ok {
					continue
				}

				identX, ok := selectorExpr.X.(*dst.Ident)
				if !ok {
					continue
				}

				if identX.Name != target.GetImportName() {
					continue
				}

				if selectorExpr.Sel.Name != target.GetNewFunction() {
					continue
				}

				if firstNewMethod == nil {
					firstNewMethod = &i
				}

				for i, dependency := range dependencies {
					if dependency == nil {
						if len(callExpr.Args) <= i {
							callExpr.Args = append(callExpr.Args,
								&dst.Ident{
									Name: "nil",
									Decs: dst.IdentDecorations{
										NodeDecs: dst.NodeDecs{
											Before: dst.NewLine,
											After:  dst.NewLine,
										},
									},
								},
							)
						}
						continue
					}

					var newIdent dst.Ident

					if _, ok := dependency.(*analyzer.AnalyzedDatabaseDependency); ok {
						newIdent = dst.Ident{
							Name: "db",
							Decs: dst.IdentDecorations{
								NodeDecs: dst.NodeDecs{
									Before: dst.NewLine,
									After:  dst.NewLine,
								},
							},
						}
					} else {
						newIdent = dst.Ident{
							Name: helpers.LowerFirstLetter(dependency.GetName()),
							Decs: dst.IdentDecorations{
								NodeDecs: dst.NodeDecs{
									Before: dst.NewLine,
									After:  dst.NewLine,
								},
							},
						}
					}

					if len(callExpr.Args) <= i {
						callExpr.Args = append(callExpr.Args, &newIdent)
					} else {
						ident, ok := callExpr.Args[i].(*dst.Ident)
						if !ok {
							continue
						}

						if ident.Name == "nil" {
							callExpr.Args[i] = &newIdent
						}
					}
				}

				break
			}
		}

		if firstNewMethod == nil {
			continue
		}

		for _, dependency := range dependencies {
			if dependency == nil {
				continue
			}

			if _, ok := dependency.(*analyzer.AnalyzedDatabaseDependency); ok {
				continue
			}

			skip := false

			for _, stmt := range funcDecl.Body.List {
				assignStmt, ok := stmt.(*dst.AssignStmt)
				if !ok {
					continue
				}

				for _, expr := range assignStmt.Rhs {
					callExpr, ok := expr.(*dst.CallExpr)
					if !ok {
						continue
					}

					selectorExpr, ok := callExpr.Fun.(*dst.SelectorExpr)
					if !ok {
						continue
					}

					identX, ok := selectorExpr.X.(*dst.Ident)
					if !ok {
						continue
					}

					if identX.Name != dependency.GetImportName() {
						continue
					}

					if selectorExpr.Sel.Name != dependency.GetNewFunction() {
						continue
					}

					skip = true

					break
				}
			}

			if skip {
				continue
			}

			insertPos := *firstNewMethod

			funcDecl.Body.List = append(funcDecl.Body.List[:insertPos], append([]dst.Stmt{
				&dst.AssignStmt{
					Lhs: []dst.Expr{
						dst.NewIdent(helpers.LowerFirstLetter(dependency.GetName())),
					},
					Rhs: []dst.Expr{
						&dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X:   dst.NewIdent(dependency.GetImportName()),
								Sel: dst.NewIdent(dependency.GetNewFunction()),
							},
						},
					},
					Tok: token.DEFINE,
				},
			}, funcDecl.Body.List[insertPos:]...)...)

			*firstNewMethod += 1
		}
	}

	if !foundFunction {
		logger.InjectorLogger.Fatal("Failed to find function `NewContainer` in container.go")
	}

	err = writeDSTFileToPath(f, path)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	for _, dependency := range dependencies {
		if dependency == nil {
			continue
		}

		InjectContainerDependencies(path, dependency)
	}
}
