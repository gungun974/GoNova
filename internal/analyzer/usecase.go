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

type AnalyzedUsecase struct {
	Name     string
	FilePath string
	PkgPath  string

	Dependencies []AnalyzedDependency
}

func (a *AnalyzedUsecase) GetDependencies() []AnalyzedDependency {
	return a.Dependencies
}

func (a *AnalyzedUsecase) GetName() string {
	return a.Name
}

func (a *AnalyzedUsecase) GetNewFunction() string {
	return "New" + helpers.CapitalizeFirstLetter(a.Name)
}

func (a *AnalyzedUsecase) GetPkgPath() string {
	return a.PkgPath
}

func (a *AnalyzedUsecase) GetImportName() string {
	return helpers.ToSnakeCase(a.Name)
}

func (a *AnalyzedUsecase) GetType() string {
	return "usecases"
}

func AnalyzeProjectUsecases(repositories []AnalyzedRepository, presenters []AnalyzedPresenter) []AnalyzedUsecase {
	usecases := []AnalyzedUsecase{}

	pkgs, err := decorator.Load(
		&packages.Config{
			Dir:  ".",
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		},
		"./internal/layers/domain/usecases/...",
	)
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	for _, pkg := range pkgs {
		rootPkgPath := pkg.PkgPath

		i := strings.LastIndex(rootPkgPath, "/")
		if i != -1 {
			rootPkgPath = rootPkgPath[:i]
		}

		if !strings.HasSuffix(rootPkgPath, "internal/layers/domain/usecases") {
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

			if !strings.HasSuffix(ident.Name, "Usecase") {
				continue
			}

			usecase := AnalyzedUsecase{
				Name:         ident.Name,
				FilePath:     pkg.Fset.Position(obj.Pos()).Filename,
				PkgPath:      pkg.PkgPath,
				Dependencies: []AnalyzedDependency{},
			}

			usecases = append(usecases, usecase)
		}
	}

	slices.SortFunc(usecases, func(a, b AnalyzedUsecase) int {
		return cmp.Compare(a.Name, b.Name)
	})

	for i := range usecases {
		DeepAnalyzeProjectUsecase(&usecases[i], repositories, presenters)
	}

	return usecases
}

func DeepAnalyzeProjectUsecase(usecase *AnalyzedUsecase, repositories []AnalyzedRepository, presenters []AnalyzedPresenter) {
	f, err := decorator.ParseFile(token.NewFileSet(), usecase.FilePath, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	usecase.Dependencies = []AnalyzedDependency{}

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

			if typeSpec.Name.Name != usecase.Name {
				continue
			}

			structType, ok := typeSpec.Type.(*dst.StructType)
			if !ok {
				continue
			}

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					usecase.Dependencies = append(usecase.Dependencies, nil)
					continue
				}

				selectorExpr, ok := field.Type.(*dst.SelectorExpr)
				if !ok {
					usecase.Dependencies = append(usecase.Dependencies, nil)
					continue
				}

				identX, ok := selectorExpr.X.(*dst.Ident)
				if !ok {
					usecase.Dependencies = append(usecase.Dependencies, nil)
					continue
				}

				found := false

				if strings.Contains(identX.Name, "repositories") {
					for _, repository := range repositories {
						if repository.Name == selectorExpr.Sel.Name {
							usecase.Dependencies = append(usecase.Dependencies, &repository)
							found = true
							break
						}
					}
				} else if strings.Contains(identX.Name, "presenters") {
					for _, presenter := range presenters {
						if presenter.Name == selectorExpr.Sel.Name {
							usecase.Dependencies = append(usecase.Dependencies, &presenter)
							found = true
							break
						}
					}
				}

				if !found {
					usecase.Dependencies = append(usecase.Dependencies, nil)
				}
			}
		}
	}
}
