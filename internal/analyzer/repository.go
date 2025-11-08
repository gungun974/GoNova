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

type AnalyzedRepository struct {
	Name     string
	FilePath string
	PkgPath  string

	Dependencies []AnalyzedDependency
}

func (a *AnalyzedRepository) GetDependencies() []AnalyzedDependency {
	return a.Dependencies
}

func (a *AnalyzedRepository) GetName() string {
	return a.Name
}

func (a *AnalyzedRepository) GetNewFunction() string {
	return "New" + helpers.CapitalizeFirstLetter(a.Name)
}

func (a *AnalyzedRepository) GetPkgPath() string {
	return a.PkgPath
}

func (a *AnalyzedRepository) GetImportName() string {
	return "repositories"
}

func (a *AnalyzedRepository) GetType() string {
	return "repositories"
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
				Name:         ident.Name,
				FilePath:     pkg.Fset.Position(obj.Pos()).Filename,
				PkgPath:      pkg.PkgPath,
				Dependencies: []AnalyzedDependency{},
			}

			repositories = append(repositories, repository)
		}
	}

	slices.SortFunc(repositories, func(a, b AnalyzedRepository) int {
		return cmp.Compare(a.Name, b.Name)
	})

	for i := range repositories {
		DeepAnalyzeProjectRepository(&repositories[i])
	}

	return repositories
}

func DeepAnalyzeProjectRepository(repository *AnalyzedRepository) {
	f, err := decorator.ParseFile(token.NewFileSet(), repository.FilePath, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	repository.Dependencies = []AnalyzedDependency{}

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

			if typeSpec.Name.Name != repository.Name {
				continue
			}

			structType, ok := typeSpec.Type.(*dst.StructType)
			if !ok {
				continue
			}

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					repository.Dependencies = append(repository.Dependencies, nil)
					continue
				}

				expr := field.Type

				starExpr, ok := field.Type.(*dst.StarExpr)
				if ok {
					expr = starExpr.X
				}

				selectorExpr, ok := expr.(*dst.SelectorExpr)
				if !ok {
					repository.Dependencies = append(repository.Dependencies, nil)
					continue
				}

				identX, ok := selectorExpr.X.(*dst.Ident)
				if !ok {
					repository.Dependencies = append(repository.Dependencies, nil)
					continue
				}

				if identX.Name == "sqlx" && selectorExpr.Sel.Name == "DB" {
					repository.Dependencies = append(repository.Dependencies, &AnalyzedDatabaseDependency{})
				} else {
					repository.Dependencies = append(repository.Dependencies, nil)
				}
			}
		}
	}
}
