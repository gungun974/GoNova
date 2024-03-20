package actions

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	core_template "github.com/gungun974/gonova/resources/core"
	server_template "github.com/gungun974/gonova/resources/core/cmd/server"
	context_template "github.com/gungun974/gonova/resources/core/internalcore/context"
	entities_template "github.com/gungun974/gonova/resources/core/internalcore/entities"
	logger_template "github.com/gungun974/gonova/resources/core/internalcore/logger"
	middlewares_template "github.com/gungun974/gonova/resources/core/internalcore/middlewares"
	models_template "github.com/gungun974/gonova/resources/core/internalcore/models"
	routes_template "github.com/gungun974/gonova/resources/core/internalcore/routes"
	layouts_template "github.com/gungun974/gonova/resources/core/resources/views/layouts"
	errors_page_template "github.com/gungun974/gonova/resources/core/resources/views/pages/errors"
	home_template "github.com/gungun974/gonova/resources/core/resources/views/pages/home"
)

func InstallCore(rawProjectName string) error {
	projectName := strings.TrimSpace(rawProjectName)

	projectPath := "."

	if _, err := os.Stat("./go.mod"); err == nil {
		logger.MainLogger.Fatal("Can't Install Core when go module is already init")
	}

	err := utils.InitGoMod(projectName, projectPath)
	if err != nil {
		return err
	}

	projectGlobalTemplateConfig := struct {
		ProjectName string
	}{
		ProjectName: projectName,
	}

	//! /

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/Makefile"), core_template.MakefileTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/.golangci.yml"), core_template.GolangciTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/.gitignore"), core_template.GitignoreTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/.air.toml"), core_template.AirTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/.env"), core_template.EnvTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /cmd/server

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/cmd/server/main.go"), server_template.MainGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /internal/context

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/context/key.go"), context_template.KeyGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/context/request_info.go"), context_template.RequestInfoGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /internal/entities

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/entities/error.go"), entities_template.ErrorGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /internal/logger

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/logger/logger.go"), logger_template.LoggerGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /internal/middlewares

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/middlewares/conditional_logger.go"), middlewares_template.ConditionalLoggerGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/middlewares/request_info_middleware.go"), middlewares_template.RequestInfoMiddlewareGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /internal/models

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/models/api.go"), models_template.ApiGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/models/html.go"), models_template.HtmlGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/models/image.go"), models_template.ImageGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/models/json.go"), models_template.JsonGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/models/pdf.go"), models_template.PdfGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/models/redirect.go"), models_template.RedirectGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /internal/routes

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/routes/file.go"), routes_template.FileGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/routes/handle_http_error.go"), routes_template.HandleHttpErrorGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/routes/home.go"), routes_template.HomeGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/internal/routes/routes.go"), routes_template.RoutesGoTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /resources/view/layouts

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/resources/views/layouts/app_layout.templ"), layouts_template.AppLayoutTemplTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/resources/views/layouts/core.templ"), layouts_template.CoreTemplTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /resources/views/pages/errors

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/resources/views/pages/errors/errors.templ"), errors_page_template.ErrorsTemplTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /resources/views/pages/home

	err = utils.CreateFileFromTemplate(filepath.Join(projectPath, "/resources/views/pages/home/index.templ"), home_template.IndexTemplTemplate, projectGlobalTemplateConfig)
	if err != nil {
		return err
	}

	//! /resources/views/pages/home

	err = utils.CreateDirectory(filepath.Join(projectPath, "/resources/views/components"))
	if err != nil {
		return err
	}

	err = utils.GoFmt(projectPath)
	if err != nil {
		return err
	}

	err = utils.GoGetPackage(projectPath, "github.com/a-h/templ/cmd/templ")
	if err != nil {
		return err
	}

	err = utils.TemplGenerate(projectPath)
	if err != nil {
		return err
	}

	err = utils.GoTidy(projectPath)
	if err != nil {
		return err
	}

	return nil
}
