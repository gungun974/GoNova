package injector

import (
	"bytes"
	"fmt"
	"go/printer"
	"go/token"
	"io"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/gungun974/gonova/internal/logger"
)

func containsImport(imp *dst.GenDecl, path string, name string) bool {
	for _, spec := range imp.Specs {
		if i, ok := spec.(*dst.ImportSpec); !ok || i == nil {
			continue
		} else if i.Path != nil && i.Path.Value == fmt.Sprintf("%q", path) {
			if i.Name != nil && i.Name.Name != name {
				continue
			}
			return true
		}
	}
	return false
}

func addImport(f *dst.File, path string, name string) bool {
	var latestImportDel *dst.GenDecl
	for _, decl := range f.Decls {
		if gen, ok := decl.(*dst.GenDecl); ok && gen != nil && gen.Tok == token.IMPORT {
			latestImportDel = gen
			if containsImport(gen, path, name) {
				return true
			}
		}
	}
	if latestImportDel == nil {
		latestImportDel = &dst.GenDecl{
			Tok:   token.IMPORT,
			Specs: []dst.Spec{},
		}
		f.Decls = append([]dst.Decl{latestImportDel}, f.Decls...)
	}
	latestImportDel.Specs = append(latestImportDel.Specs, &dst.ImportSpec{
		Name: dst.NewIdent(name),
		Path: &dst.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("%q", path),
		},
	})
	return false
}

func generateDSTFileContent(file *dst.File) (string, error) {
	var buf bytes.Buffer
	writer := io.Writer(&buf)

	fset, af, err := decorator.RestoreFile(file)
	if err != nil {
		return "", err
	}

	if err := printer.Fprint(writer, fset, af); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func writeDSTFileToPath(file *dst.File, path string) error {
	content, err := generateDSTFileContent(file)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, []byte(content), 0)
	if err != nil {
		logger.MainLogger.Errorf("Failed to write in file %s : %v", path, err)
		return err
	}

	return nil
}
