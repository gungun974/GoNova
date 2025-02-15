package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	make_usecase_template "github.com/gungun974/gonova/resources/make/usecase"
)

func MakeUsecase(usecaseName string) error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	newUsecaseFilePath := fmt.Sprintf(
		"/internal/layers/domain/usecases/%s/%s_usecase.go",
		helpers.ToSnakeCase(usecaseName),
		helpers.ToSnakeCase(usecaseName),
	)

	usecaseExampleFilePath := fmt.Sprintf(
		"/internal/layers/domain/usecases/%s/%s_get_example.go",
		helpers.ToSnakeCase(usecaseName),
		helpers.ToSnakeCase(usecaseName),
	)

	if _, err := os.Stat(filepath.Join(projectPath, newUsecaseFilePath)); err == nil {
		logger.MainLogger.Fatal(
			"Can't Make Usecase when a usecase with the same name is already present",
		)
	}

	logger.MainLogger.Info("Make Usecase")

	projectGlobalTemplateConfig := struct {
		ProjectName    string
		UsecaseName    string
		UsecasePackage string
	}{
		ProjectName:    projectName,
		UsecaseName:    helpers.CapitalizeFirstLetter(usecaseName) + "Usecase",
		UsecasePackage: strings.ToLower(usecaseName) + "_usecase",
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, newUsecaseFilePath),
		make_usecase_template.UsecaseGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, usecaseExampleFilePath),
		make_usecase_template.UsecaseGetExampleGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.GoFmt(projectPath)
	if err != nil {
		return err
	}

	return nil
}
