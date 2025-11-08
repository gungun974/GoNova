package actions

import (
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/utils"
)

func LinkUsecase(usecase analyzer.AnalyzedUsecase, controller analyzer.AnalyzedController) error {
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

	return nil
}
