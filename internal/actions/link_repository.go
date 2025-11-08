package actions

import (
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/utils"
)

func LinkRepository(repository analyzer.AnalyzedRepository, usecase analyzer.AnalyzedUsecase) error {
	usecaseFilePath := usecase.FilePath

	injector.InjectUsecaseRepository(usecaseFilePath, usecase, repository)

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
