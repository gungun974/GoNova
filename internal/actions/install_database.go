package actions

import (
	"os"
	"path/filepath"

	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	database_template "github.com/gungun974/gonova/resources/database"
)

func InstallPostgreDatabase() error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	if _, err := os.Stat(filepath.Join(projectPath, "/internal/database/database.go")); err == nil {
		logger.MainLogger.Fatal(
			"Can't Install Database Module when internal database.go is already installed",
		)
	}

	if !utils.VerifyCmd("go") {
		logger.MainLogger.Fatal("Can't Install Postgre Database without go in PATH")
	}

	logger.MainLogger.Info("Install Postgre Database")

	projectGlobalTemplateConfig := struct {
		ProjectName   string
		EnablePostgre bool
		EnableSqlite  bool
	}{
		ProjectName:   projectName,
		EnablePostgre: true,
	}

	err = utils.MergeFileFromTemplate(
		filepath.Join(projectPath, "/internal/logger/logger.go"),
		database_template.LoggerGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/database/database.go"),
		database_template.PostgreDatabaseGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/database/migrations/EMPTY_MIGRATION.sql"),
		database_template.EmptyMigrationTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.MergeFileFromTemplate(
		filepath.Join(projectPath, "/.env"),
		database_template.PostgreEnvFileTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	dockerComposeContent, err := utils.TemplateToString(
		database_template.PostgreDockerComposeFileTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.MergeYAMLFromString(
		filepath.Join(projectPath, "/docker-compose.yml"),
		dockerComposeContent,
	)
	if err != nil {
		return err
	}

	injector.InjectMainDatabase(
		filepath.Join(projectPath, "/cmd/server/main.go"),
		projectName,
	)

	injector.InjectContainerDatabase(
		filepath.Join(projectPath, "/internal/container.go"),
	)

	err = utils.GoImports(projectPath)
	if err != nil {
		return err
	}

	err = utils.GoFumpt(projectPath)
	if err != nil {
		return err
	}

	return nil
}

func InstallSqliteDatabase() error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	if _, err := os.Stat(filepath.Join(projectPath, "/internal/database/database.go")); err == nil {
		logger.MainLogger.Fatal(
			"Can't Install Database Module when internal database.go is already installed",
		)
	}

	if !utils.VerifyCmd("go") {
		logger.MainLogger.Fatal("Can't Install Sqlite Database without go in PATH")
	}

	logger.MainLogger.Info("Install Sqlite Database")

	projectGlobalTemplateConfig := struct {
		ProjectName   string
		EnablePostgre bool
		EnableSqlite  bool
	}{
		ProjectName:  projectName,
		EnableSqlite: true,
	}

	err = utils.MergeFileFromTemplate(
		filepath.Join(projectPath, "/internal/logger/logger.go"),
		database_template.LoggerGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/database/database.go"),
		database_template.SqliteDatabaseGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/database/migrations/EMPTY_MIGRATION.sql"),
		database_template.EmptyMigrationTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.MergeFileFromTemplate(
		filepath.Join(projectPath, "/.env"),
		database_template.SqliteEnvTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateDirectory(filepath.Join(projectPath, "/data"))
	if err != nil {
		return err
	}

	injector.InjectMainDatabase(
		filepath.Join(projectPath, "/cmd/server/main.go"),
		projectName,
	)

	injector.InjectContainerDatabase(
		filepath.Join(projectPath, "/internal/container.go"),
	)

	err = utils.GoImports(projectPath)
	if err != nil {
		return err
	}

	err = utils.GoFumpt(projectPath)
	if err != nil {
		return err
	}

	return nil
}
