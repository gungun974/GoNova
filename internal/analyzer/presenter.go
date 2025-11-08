package analyzer

import (
	"cmp"
	"go/types"
	"slices"
	"strings"

	"github.com/dave/dst/decorator"
	"github.com/gungun974/gonova/internal/logger"
	"golang.org/x/tools/go/packages"
)

type AnalyzedPresenter struct {
	Name     string
	FilePath string
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
				Name:     ident.Name,
				FilePath: pkg.Fset.Position(obj.Pos()).Filename,
			}

			presenters = append(presenters, presenter)
		}
	}

	slices.SortFunc(presenters, func(a, b AnalyzedPresenter) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return presenters
}
