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

const APP_ERRROR_STRUCT_NAME string = "AppError"

type AnalyzedEntity struct {
	Name     string
	FilePath string
	Fields   []*types.Var
}

func (a AnalyzedEntity) Equal(b AnalyzedEntity) bool {
	return a.Name == b.Name && a.FilePath == b.FilePath
}

func AnalyzeProjectEntities() []AnalyzedEntity {
	entities := []AnalyzedEntity{}

	pkgs, err := decorator.Load(
		&packages.Config{
			Dir:  ".",
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		},
		"./internal/layers/domain/entities/...",
	)
	if err != nil {
		logger.MainLogger.Logger.Fatal(err)
	}

	var appErrorType *types.Named

	for _, pkg := range pkgs {
		if !strings.HasSuffix(pkg.PkgPath, "internal/layers/domain/entities") {
			continue
		}

		for ident, obj := range pkg.TypesInfo.Defs {
			if obj == nil {
				continue
			}
			if typeName, ok := obj.(*types.TypeName); ok {
				if namedType, ok := typeName.Type().(*types.Named); ok {
					if ident.Name == APP_ERRROR_STRUCT_NAME {
						appErrorType = namedType
					}
				}
			}
		}
	}

	for _, pkg := range pkgs {
		if !strings.HasSuffix(pkg.PkgPath, "internal/layers/domain/entities") {
			continue
		}

	defs_loop:
		for ident, obj := range pkg.TypesInfo.Defs {
			if obj == nil {
				continue
			}

			typeName, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}
			if namedType, ok := typeName.Type().(*types.Named); ok {
				if namedType.Underlying() == appErrorType.Underlying() {
					continue
				}
			}

			if structType, ok := typeName.Type().Underlying().(*types.Struct); ok {
				for i := 0; i < structType.NumFields(); i++ {
					field := structType.Field(i)
					if field.Anonymous() && field.Type() == appErrorType {
						continue defs_loop
					}
				}
			}

			if !typeName.Exported() {
				continue
			}

			entity := AnalyzedEntity{
				Name:     ident.Name,
				FilePath: pkg.Fset.Position(obj.Pos()).Filename,
			}

			if structType, ok := typeName.Type().Underlying().(*types.Struct); ok {
				for i := 0; i < structType.NumFields(); i++ {
					field := structType.Field(i)
					entity.Fields = append(entity.Fields, field)
				}
			}

			entities = append(entities, entity)
		}
	}

	slices.SortFunc(entities, func(a, b AnalyzedEntity) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return entities
}
