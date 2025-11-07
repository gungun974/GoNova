package actions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	make_controller_template "github.com/gungun974/gonova/resources/make/controller"
)

func MakeController(controllerName string) error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	newControllerFilePath := fmt.Sprintf(
		"/internal/layers/presentation/controllers/%s_controller.go",
		helpers.ToSnakeCase(controllerName),
	)

	if _, err := os.Stat(filepath.Join(projectPath, newControllerFilePath)); err == nil {
		logger.MainLogger.Fatal(
			"Can't Make Controller when a controller with the same name is already present",
		)
	}

	logger.MainLogger.Info("Make Controller")

	projectGlobalTemplateConfig := struct {
		ProjectName    string
		ControllerName string
	}{
		ProjectName:    projectName,
		ControllerName: helpers.CapitalizeFirstLetter(controllerName) + "Controller",
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, newControllerFilePath),
		make_controller_template.ControllerGoTemplate,
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
