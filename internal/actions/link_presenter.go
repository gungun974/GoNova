package actions

import (
	"path/filepath"

	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
)

func LinkPresenter(presenter analyzer.AnalyzedPresenter, usecase analyzer.AnalyzedUsecase) error {
	projectPath := "."

	containerFilePath := filepath.Join(projectPath, "/internal/container.go")

	usecaseFilePath := usecase.FilePath

	logger.MainLogger.Info("Link Presenter")

	injector.InjectUsecasePresenter(usecaseFilePath, usecase, presenter)

	err := utils.GoImports(usecaseFilePath)
	if err != nil {
		return err
	}

	err = utils.GoFumpt(usecaseFilePath)
	if err != nil {
		return err
	}

	repositories := analyzer.AnalyzeProjectRepositories()
	presenters := analyzer.AnalyzeProjectPresenters()
	analyzer.DeepAnalyzeProjectUsecase(&usecase, repositories, presenters)

	injector.InjectContainerDependencies(containerFilePath, &usecase)

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
