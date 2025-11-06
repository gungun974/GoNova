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

type AnalyzedModel struct {
	Name     string
	FilePath string
	Entity   AnalyzedEntity
}

func AnalyzeProjectModels(entities []AnalyzedEntity) []AnalyzedModel {
	models := []AnalyzedModel{}

	pkgs, err := decorator.Load(
		&packages.Config{
			Dir:  ".",
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		},
		"./internal/layers/data/models/...",
	)
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	for _, pkg := range pkgs {
		if !strings.HasSuffix(pkg.PkgPath, "internal/layers/data/models") {
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

			if !strings.HasSuffix(ident.Name, "Model") {
				continue
			}

			model := AnalyzedModel{
				Name:     ident.Name,
				FilePath: pkg.Fset.Position(obj.Pos()).Filename,
			}

			for _, entity := range entities {
				if entity.Name != strings.TrimSuffix(model.Name, "Model") {
					continue
				}
				model.Entity = entity
				break
			}

			models = append(models, model)
		}
	}

	slices.SortFunc(models, func(a, b AnalyzedModel) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return models
}
