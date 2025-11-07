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

type AnalyzedRepository struct {
	Name     string
	FilePath string
}

func AnalyzeProjectRepositories() []AnalyzedRepository {
	repositories := []AnalyzedRepository{}

	pkgs, err := decorator.Load(
		&packages.Config{
			Dir:  ".",
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		},
		"./internal/layers/data/repositories/...",
	)
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	for _, pkg := range pkgs {
		if !strings.HasSuffix(pkg.PkgPath, "internal/layers/data/repositories") {
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

			if !strings.HasSuffix(ident.Name, "Repository") {
				continue
			}

			repository := AnalyzedRepository{
				Name:     ident.Name,
				FilePath: pkg.Fset.Position(obj.Pos()).Filename,
			}

			repositories = append(repositories, repository)
		}
	}

	slices.SortFunc(repositories, func(a, b AnalyzedRepository) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return repositories
}
