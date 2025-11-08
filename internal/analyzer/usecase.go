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

type AnalyzedUsecase struct {
	Name     string
	FilePath string
	PkgPath  string
}

func AnalyzeProjectUsecases() []AnalyzedUsecase {
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
				Name:     ident.Name,
				FilePath: pkg.Fset.Position(obj.Pos()).Filename,
				PkgPath:  pkg.PkgPath,
			}

			usecases = append(usecases, usecase)
		}
	}

	slices.SortFunc(usecases, func(a, b AnalyzedUsecase) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return usecases
}
