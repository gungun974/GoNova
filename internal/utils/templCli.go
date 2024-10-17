package utils

import (
	"github.com/gungun974/gonova/internal/logger"
)

func TemplGenerate(appDir string) error {
	templVersion, err := GetTemplVersion(appDir)
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	if err := ExecuteCmd("go",
		[]string{"run", "github.com/a-h/templ/cmd/templ@" + templVersion, "generate", "-v"},
		appDir); err != nil {
		return err
	}

	return nil
}
