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

type AnalyzedPresenter struct {
	Name     string
	FilePath string
	PkgPath  string

	Dependencies []AnalyzedDependency
}

func (a *AnalyzedPresenter) GetDependencies() []AnalyzedDependency {
	return a.Dependencies
}

func (a *AnalyzedPresenter) GetName() string {
	return a.Name
}

func (a *AnalyzedPresenter) GetNewFunction() string {
	return "New" + helpers.CapitalizeFirstLetter(a.Name)
}

func (a *AnalyzedPresenter) GetPkgPath() string {
	return a.PkgPath
}

func (a *AnalyzedPresenter) GetImportName() string {
	return "presenters"
}

func AnalyzeProjectPresenters() []AnalyzedPresenter {
	presenters := []AnalyzedPresenter{}

	pkgs, err := decorator.Load(
		&packages.Config{
			Dir:  ".",
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		},
		"./internal/layers/presentation/presenters/...",
	)
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	for _, pkg := range pkgs {
		if !strings.HasSuffix(pkg.PkgPath, "internal/layers/presentation/presenters") {
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

			if !strings.HasSuffix(ident.Name, "Presenter") {
				continue
			}

			presenter := AnalyzedPresenter{
				Name:         ident.Name,
				FilePath:     pkg.Fset.Position(obj.Pos()).Filename,
				PkgPath:      pkg.PkgPath,
				Dependencies: []AnalyzedDependency{},
			}

			f, err := decorator.ParseFile(token.NewFileSet(), presenter.FilePath, nil, parser.ParseComments)
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

					if typeSpec.Name.Name != presenter.Name {
						continue
					}

					structType, ok := typeSpec.Type.(*dst.StructType)
					if !ok {
						continue
					}

					for range structType.Fields.List {
						presenter.Dependencies = append(presenter.Dependencies, nil)
					}
				}
			}

			presenters = append(presenters, presenter)
		}
	}

	slices.SortFunc(presenters, func(a, b AnalyzedPresenter) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return presenters
}
