package actions

import (
	"path/filepath"

	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/utils"
)

func LinkUsecase(usecase analyzer.AnalyzedUsecase, controller analyzer.AnalyzedController) error {
	projectPath := "."

	containerFilePath := filepath.Join(projectPath, "/internal/container.go")

	controllerFilePath := controller.FilePath

	injector.InjectControllerUsecase(controllerFilePath, controller, usecase)

	err := utils.GoImports(controllerFilePath)
	if err != nil {
		return err
	}

	err = utils.GoFumpt(controllerFilePath)
	if err != nil {
		return err
	}

	repositories := analyzer.AnalyzeProjectRepositories()
	presenters := analyzer.AnalyzeProjectPresenters()
	usecases := analyzer.AnalyzeProjectUsecases(repositories, presenters)
	analyzer.DeepAnalyzeProjectController(&controller, usecases)

	injector.InjectContainerDependencies(containerFilePath, &controller)

	err = utils.GoImports(containerFilePath)
	if err != nil {
		return err
	}

	err = utils.GoFumpt(containerFilePath)
	if err != nil {
		return err
	}

	return nil
}
