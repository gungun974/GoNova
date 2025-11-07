package actions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	make_presenter_template "github.com/gungun974/gonova/resources/make/presenter"
)

func MakePresenter(presenterName string) error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	newPresenterFilePath := fmt.Sprintf(
		"/internal/layers/presentation/presenters/%s_presenter.go",
		helpers.ToSnakeCase(presenterName),
	)

	if _, err := os.Stat(filepath.Join(projectPath, newPresenterFilePath)); err == nil {
		logger.MainLogger.Fatal(
			"Can't Make Presenter when a presenter with the same name is already present",
		)
	}

	logger.MainLogger.Info("Make Presenter")

	projectGlobalTemplateConfig := struct {
		ProjectName   string
		PresenterName string
	}{
		ProjectName:   projectName,
		PresenterName: helpers.CapitalizeFirstLetter(presenterName) + "Presenter",
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, newPresenterFilePath),
		make_presenter_template.PresenterGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.GoFumpt(projectPath)
	if err != nil {
		return err
	}

	return nil
}
