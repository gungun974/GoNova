package actions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	make_model_template "github.com/gungun974/gonova/resources/make/model"

	"github.com/jinzhu/inflection"
)

func MakeModel(entity analyzer.AnalyzedEntity) error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	newEntityFilePath := fmt.Sprintf(
		"/internal/layers/data/models/%s.go",
		helpers.ToSnakeCase(entity.Name),
	)

	logger.MainLogger.Info("Make Model")

	projectGlobalTemplateConfig := struct {
		ProjectName string

		ModelName  string
		ModelsName string

		EntityName   string
		EntitiesName string
	}{
		ProjectName: projectName,

		ModelName:  helpers.CapitalizeFirstLetter(entity.Name) + "Model",
		ModelsName: inflection.Plural(helpers.CapitalizeFirstLetter(entity.Name) + "Model"),

		EntityName:   helpers.CapitalizeFirstLetter(entity.Name),
		EntitiesName: inflection.Plural(helpers.CapitalizeFirstLetter(entity.Name)),
	}

	if _, err := os.Stat(filepath.Join(projectPath, newEntityFilePath)); err != nil {
		err = utils.CreateFileFromTemplate(
			filepath.Join(projectPath, newEntityFilePath),
			make_model_template.ModelGoTemplate,
			projectGlobalTemplateConfig,
		)
		if err != nil {
			return err
		}
	}

	err = utils.GoFmt(projectPath)
	if err != nil {
		return err
	}

	return nil
}
