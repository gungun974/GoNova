package actions

import (
	"path/filepath"

	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/utils"
)

func LinkController(repositoryName string, controller analyzer.AnalyzedController) error {
	projectPath := "."

	containerFilePath := filepath.Join(projectPath, "/internal/container.go")

	injector.InjectContainerController(containerFilePath, controller)

	err := utils.GoImports(filepath.Join(projectPath, containerFilePath))
	if err != nil {
		return err
	}

	err = utils.GoFumpt(filepath.Join(projectPath, containerFilePath))
	if err != nil {
		return err
	}

	return nil
}
