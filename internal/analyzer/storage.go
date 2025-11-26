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

type AnalyzedStorage struct {
	Name     string
	FilePath string
	PkgPath  string

	Dependencies []AnalyzedDependency
}

func (a *AnalyzedStorage) GetDependencies() []AnalyzedDependency {
	return a.Dependencies
}

func (a *AnalyzedStorage) GetName() string {
	return a.Name
}

func (a *AnalyzedStorage) GetNewFunction() string {
	return "New" + helpers.CapitalizeFirstLetter(a.Name)
}

func (a *AnalyzedStorage) GetPkgPath() string {
	return a.PkgPath
}

func (a *AnalyzedStorage) GetImportName() string {
	return "storages"
}

func (a *AnalyzedStorage) GetType() string {
	return "storages"
}

func AnalyzeProjectStorages() []AnalyzedStorage {
	storages := []AnalyzedStorage{}

	pkgs, err := decorator.Load(
		&packages.Config{
			Dir:  ".",
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		},
		"./internal/layers/data/storages/...",
	)
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	for _, pkg := range pkgs {
		if !strings.HasSuffix(pkg.PkgPath, "internal/layers/data/storages") {
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

			if !strings.HasSuffix(ident.Name, "Storage") {
				continue
			}

			storage := AnalyzedStorage{
				Name:         ident.Name,
				FilePath:     pkg.Fset.Position(obj.Pos()).Filename,
				PkgPath:      pkg.PkgPath,
				Dependencies: []AnalyzedDependency{},
			}

			storages = append(storages, storage)
		}
	}

	slices.SortFunc(storages, func(a, b AnalyzedStorage) int {
		return cmp.Compare(a.Name, b.Name)
	})

	for i := range storages {
		DeepAnalyzeProjectStorage(&storages[i])
	}

	return storages
}

func DeepAnalyzeProjectStorage(storage *AnalyzedStorage) {
	f, err := decorator.ParseFile(token.NewFileSet(), storage.FilePath, nil, parser.ParseComments)
	if err != nil {
		logger.InjectorLogger.Fatal(err)
	}

	storage.Dependencies = []AnalyzedDependency{}

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

			if typeSpec.Name.Name != storage.Name {
				continue
			}

			structType, ok := typeSpec.Type.(*dst.StructType)
			if !ok {
				continue
			}

			for range structType.Fields.List {
				storage.Dependencies = append(storage.Dependencies, nil)
			}
		}
	}
}
