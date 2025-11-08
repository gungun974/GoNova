package analyzer

import (
	"cmp"
	"go/parser"
	"go/token"
	"go/types"
	"slices"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"golang.org/x/tools/go/packages"
)

type AnalyzedController struct {
	Name     string
	FilePath string
	PkgPath  string

	Dependencies []AnalyzedDependency
}

func (a *AnalyzedController) GetDependencies() []AnalyzedDependency {
	return a.Dependencies
}

func (a *AnalyzedController) GetName() string {
	return a.Name
}

func (a *AnalyzedController) GetNewFunction() string {
	return "New" + helpers.CapitalizeFirstLetter(a.Name)
}

func (a *AnalyzedController) GetPkgPath() string {
	return a.PkgPath
}

func (a *AnalyzedController) GetImportName() string {
	return "controllers"
}

func AnalyzeProjectControllers(usecases []AnalyzedUsecase) []AnalyzedController {
	controllers := []AnalyzedController{}

	pkgs, err := decorator.Load(
		&packages.Config{
			Dir:  ".",
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		},
		"./internal/layers/presentation/controllers/...",
	)
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	for _, pkg := range pkgs {
		if !strings.HasSuffix(pkg.PkgPath, "internal/layers/presentation/controllers") {
			continue
		}

		for ident, obj := range pkg.TypesInfo.Defs {
			if obj == nil {
				continue
			}

			typeName, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}

			if !typeName.Exported() {
				continue
			}

			if !strings.HasSuffix(ident.Name, "Controller") {
				continue
			}

			controller := AnalyzedController{
				Name:         ident.Name,
				FilePath:     pkg.Fset.Position(obj.Pos()).Filename,
				PkgPath:      pkg.PkgPath,
				Dependencies: []AnalyzedDependency{},
			}

			f, err := decorator.ParseFile(token.NewFileSet(), controller.FilePath, nil, parser.ParseComments)
			if err != nil {
				logger.InjectorLogger.Fatal(err)
			}

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

					if typeSpec.Name.Name != controller.Name {
						continue
					}

					structType, ok := typeSpec.Type.(*dst.StructType)
					if !ok {
						continue
					}

					for _, field := range structType.Fields.List {
						if len(field.Names) == 0 {
							controller.Dependencies = append(controller.Dependencies, nil)
							continue
						}

						selectorExpr, ok := field.Type.(*dst.SelectorExpr)
						if !ok {
							controller.Dependencies = append(controller.Dependencies, nil)
							continue
						}

						identX, ok := selectorExpr.X.(*dst.Ident)
						if !ok {
							controller.Dependencies = append(controller.Dependencies, nil)
							continue
						}

						if !strings.Contains(identX.Name, "usecase") {
							controller.Dependencies = append(controller.Dependencies, nil)
							continue
						}

						if !strings.HasSuffix(selectorExpr.Sel.Name, "Usecase") {
							controller.Dependencies = append(controller.Dependencies, nil)
							continue
						}

						found := false

						for _, usecase := range usecases {
							if usecase.Name == selectorExpr.Sel.Name {
								controller.Dependencies = append(controller.Dependencies, &usecase)
								found = true
								break
							}
						}

						if !found {
							controller.Dependencies = append(controller.Dependencies, nil)
						}
					}
				}
			}

			controllers = append(controllers, controller)
		}
	}

	slices.SortFunc(controllers, func(a, b AnalyzedController) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return controllers
}
