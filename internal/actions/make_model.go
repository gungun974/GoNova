package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	make_model_template "github.com/gungun974/gonova/resources/make/model"
)

func MakeModel(entity analyzer.AnalyzedEntity) error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	filename := filepath.Base(entity.FilePath)

	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	newModelFilePath := fmt.Sprintf(
		"/internal/layers/data/models/%s.go",
		name,
	)

	logger.MainLogger.Info("Make Model")

	projectGlobalTemplateConfig := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}

	if _, err := os.Stat(filepath.Join(projectPath, newModelFilePath)); err != nil {
		err = utils.CreateFileFromTemplate(
			filepath.Join(projectPath, newModelFilePath),
			make_model_template.BlankModelGoTemplate,
			projectGlobalTemplateConfig,
		)
		if err != nil {
			return err
		}
	}

	injector.InjectModelNewModel(
		filepath.Join(projectPath, newModelFilePath),
		entity,
	)

	err = utils.GoImports(filepath.Join(projectPath, newModelFilePath))
	if err != nil {
		return err
	}

	err = utils.GoFumpt(filepath.Join(projectPath, newModelFilePath))
	if err != nil {
		return err
	}

	return nil
}
