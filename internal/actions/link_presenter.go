package actions

import (
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/utils"
)

func LinkPresenter(presenter analyzer.AnalyzedPresenter, usecase analyzer.AnalyzedUsecase) error {
	usecaseFilePath := usecase.FilePath

	injector.InjectUsecasePresenter(usecaseFilePath, usecase, presenter)

	err := utils.GoImports(usecaseFilePath)
	if err != nil {
		return err
	}

	err = utils.GoFumpt(usecaseFilePath)
	if err != nil {
		return err
	}

	return nil
}
