package actions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	make_entity_template "github.com/gungun974/gonova/resources/make/entity"
)

func MakeEntity(entityName string) error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	newEntityFilePath := fmt.Sprintf(
		"/internal/layers/domain/entities/%s.go",
		helpers.ToSnakeCase(entityName),
	)

	logger.MainLogger.Info("Make Entity")

	projectGlobalTemplateConfig := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}

	if _, err := os.Stat(filepath.Join(projectPath, newEntityFilePath)); err != nil {
		err = utils.CreateFileFromTemplate(
			filepath.Join(projectPath, newEntityFilePath),
			make_entity_template.EntityGoTemplate,
			projectGlobalTemplateConfig,
		)
		if err != nil {
			return err
		}
	}

	injector.InjectEntityNewEntity(
		filepath.Join(projectPath, newEntityFilePath),
		helpers.CapitalizeFirstLetter(entityName),
	)

	err = utils.GoImports(filepath.Join(projectPath, newEntityFilePath))
	if err != nil {
		return err
	}

	err = utils.GoFumpt(filepath.Join(projectPath, newEntityFilePath))
	if err != nil {
		return err
	}

	return nil
}
