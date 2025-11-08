package actions

import (
	"path/filepath"

	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
)

func LinkController(controller analyzer.AnalyzedController) error {
	projectPath := "."

	containerFilePath := filepath.Join(projectPath, "/internal/container.go")

	logger.MainLogger.Info("Link Controller")

	injector.InjectContainerController(containerFilePath, controller)

	injector.InjectContainerDependencies(containerFilePath, &controller)

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
