package utils

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func GetGoModName(appDir string) (string, error) {
	goModPath := filepath.Join(appDir, "/go.mod")

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	file, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		return "", err
	}

	// Récupère et affiche le nom du module
	if file.Module != nil && file.Module.Mod.Path != "" {
		return file.Module.Mod.Path, nil
	}

	return "", errors.New("Can't find go module name")
}

func GetTemplVersion(appDir string) (string, error) {
	goModPath := filepath.Join(appDir, "/go.mod")

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	file, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		return "", err
	}

	for _, require := range file.Require {
		if require.Mod.Path == "github.com/a-h/templ" {
			return require.Mod.Version, nil
		}
	}

	return "", errors.New("Can't find templ module version")
}
