package utils

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gungun974/gonova/internal/logger"
	"golang.org/x/tools/imports"

	gofumpt "mvdan.cc/gofumpt/format"
)

func GoImports(path string) error {
	switch info, err := os.Stat(path); {
	case err != nil:
		return err
	case !info.IsDir():
		return goImportsFile(path)
	default:
		// Directories are walked, ignoring non-Go files.
		err := filepath.WalkDir(path, func(newPath string, f fs.DirEntry, err error) error {
			// vendor and testdata directories are skipped,
			// unless they are explicitly passed as an argument.
			base := filepath.Base(newPath)
			if newPath != path && (base == "vendor" || base == "testdata") {
				return filepath.SkipDir
			}

			if err != nil {
				return err
			}

			if !isGoFile(f) {
				return nil
			}
			_, err = f.Info()
			if err != nil {
				return err
			}
			return goImportsFile(newPath)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func goImportsFile(path string) error {
	formatted, err := imports.Process(path, nil, nil)
	if err != nil {
		logger.CommandLogger.Errorf("Failed to use goimports in file %s : %v", path, err)
		return err
	}

	err = os.WriteFile(path, []byte(formatted), 0)
	if err != nil {
		logger.CommandLogger.Errorf("Failed to write in file %s : %v", path, err)
		return err
	}

	return nil
}

func GoFumpt(path string) error {
	goModPath := filepath.Join(".", "/go.mod")

	switch info, err := os.Stat(path); {
	case err != nil:
		return err
	case !info.IsDir():
		return goFumptFile(goModPath, path)
	default:
		// Directories are walked, ignoring non-Go files.
		err := filepath.WalkDir(path, func(newPath string, f fs.DirEntry, err error) error {
			// vendor and testdata directories are skipped,
			// unless they are explicitly passed as an argument.
			base := filepath.Base(newPath)
			if newPath != path && (base == "vendor" || base == "testdata") {
				return filepath.SkipDir
			}

			if err != nil || !isGoFile(f) {
				return err
			}
			_, err = f.Info()
			if err != nil {
				return err
			}
			return goFumptFile(goModPath, newPath)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func isGoFile(f fs.DirEntry) bool {
	// ignore non-Go files
	name := f.Name()
	return !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go") && !f.IsDir()
}

func goFumptFile(goModPath string, path string) error {
	fset := token.NewFileSet()

	// Ensure our parsed files never start with base 1,
	// to ensure that using token.NoPos+1 will panic.
	fset.AddFile("gofumpt_base.go", 1, 10)

	file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution|parser.ParseComments)
	if err != nil {
		return err
	}

	gofumpt.File(fset, file, gofumpt.Options{ModulePath: goModPath})

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return err
	}

	err = os.WriteFile(path, buf.Bytes(), 0)
	if err != nil {
		logger.MainLogger.Errorf("Failed to write in file %s : %v", path, err)
		return err
	}

	return nil
}
