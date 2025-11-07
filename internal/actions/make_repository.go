package actions

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	make_repository_template "github.com/gungun974/gonova/resources/make/repository"
)

func MakeRepository(repositoryName string) error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	newRepositoryFilePath := fmt.Sprintf(
		"/internal/layers/presentation/repositories/%s_repository.go",
		helpers.ToSnakeCase(repositoryName),
	)

	if _, err := os.Stat(filepath.Join(projectPath, newRepositoryFilePath)); err == nil {
		logger.MainLogger.Fatal(
			"Can't Make Repository when a repository with the same name is already present",
		)
	}

	logger.MainLogger.Info("Make Repository")

	projectGlobalTemplateConfig := struct {
		ProjectName    string
		RepositoryName string
	}{
		ProjectName:    projectName,
		RepositoryName: helpers.CapitalizeFirstLetter(repositoryName) + "Repository",
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, newRepositoryFilePath),
		make_repository_template.EmptyRepositoryGoTemplate,
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
