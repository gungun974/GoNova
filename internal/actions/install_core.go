package actions

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	core_template "github.com/gungun974/gonova/resources/core"
	server_template "github.com/gungun974/gonova/resources/core/cmd/server"
	internal_template "github.com/gungun974/gonova/resources/core/internalcore"
	context_template "github.com/gungun974/gonova/resources/core/internalcore/context"
	entities_template "github.com/gungun974/gonova/resources/core/internalcore/entities"
	logger_template "github.com/gungun974/gonova/resources/core/internalcore/logger"
	middlewares_template "github.com/gungun974/gonova/resources/core/internalcore/middlewares"
	models_template "github.com/gungun974/gonova/resources/core/internalcore/models"
	routes_template "github.com/gungun974/gonova/resources/core/internalcore/routes"
	resources_template "github.com/gungun974/gonova/resources/core/resources"
	css_template "github.com/gungun974/gonova/resources/core/resources/css"
	scripts_template "github.com/gungun974/gonova/resources/core/resources/scripts"
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

	if !utils.VerifyCmd("go") {
		logger.MainLogger.Fatal("Can't Install Core without go in PATH")
	}

	if !utils.VerifyCmd("pnpm") {
		logger.MainLogger.Fatal("Can't Install Core without pnpm in PATH")
	}

	logger.MainLogger.Info("Install Core")

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

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/Makefile"),
		core_template.MakefileTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/.golangci.yml"),
		core_template.GolangciTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/.gitignore"),
		core_template.GitignoreTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/.env"),
		core_template.EnvTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateDirectory(filepath.Join(projectPath, "/public"))
	if err != nil {
		return err
	}

	// JS / CSS

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/eslint.config.mjs"),
		core_template.EslintConfigFileTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/.npmrc"),
		core_template.NpmrcTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/package.json"),
		core_template.PackageJsonTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/pnpm-workspace.yaml"),
		core_template.PnpmWorkspaceTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/prettier.config.js"),
		core_template.PrettierTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/tsconfig.json"),
		core_template.TsconfigTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/vite.config.ts"),
		core_template.ViteTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /internal

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/container.go"),
		internal_template.ContainerGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /cmd/server

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/cmd/server/main.go"),
		server_template.MainGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /internal/context

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/context/key.go"),
		context_template.KeyGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/context/request_info.go"),
		context_template.RequestInfoGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/context/vite.go"),
		context_template.ViteGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /internal/layers/domain/entities

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/layers/domain/entities/error.go"),
		entities_template.ErrorGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /internal/logger

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/logger/logger.go"),
		logger_template.LoggerGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /internal/middlewares

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/middlewares/conditional_logger.go"),
		middlewares_template.ConditionalLoggerGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/middlewares/request_info_middleware.go"),
		middlewares_template.RequestInfoMiddlewareGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/middlewares/vite.go"),
		middlewares_template.ViteGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}
	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/middlewares/vite_debug.go"),
		middlewares_template.ViteDebugGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /internal/models

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/models/api.go"),
		models_template.ApiGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/models/file.go"),
		models_template.FileGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/models/html.go"),
		models_template.HtmlGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/models/image.go"),
		models_template.ImageGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/models/json.go"),
		models_template.JsonGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/models/pdf.go"),
		models_template.PdfGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/models/plain.go"),
		models_template.PlainGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/models/reader.go"),
		models_template.ReaderGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/models/redirect.go"),
		models_template.RedirectGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /internal/routes

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/routes/file.go"),
		routes_template.FileGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/routes/handle_http_error.go"),
		routes_template.HandleHttpErrorGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/routes/home.go"),
		routes_template.HomeGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/internal/routes/routes.go"),
		routes_template.RoutesGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /resources/view/layouts

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/views/layouts/app_layout.templ"),
		layouts_template.AppLayoutTemplTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/views/layouts/core.templ"),
		layouts_template.CoreTemplTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/views/layouts/vite.templ"),
		layouts_template.ViteTemplTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /resources/views/pages/errors

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/views/pages/errors/errors.templ"),
		errors_page_template.ErrorsTemplTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /resources/views/pages/home

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/views/pages/home/index.templ"),
		home_template.IndexTemplTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /resources/views/pages/home

	err = utils.CreateDirectory(filepath.Join(projectPath, "/resources/views/components"))
	if err != nil {
		return err
	}

	//! /resources

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/dom.d.ts"),
		resources_template.DomDTSTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/main.ts"),
		resources_template.MainTSTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /resources/scripts

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/scripts/core.ts"),
		scripts_template.CoreTSTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/scripts/htmx/core.ts"),
		scripts_template.HtmxCoreTSTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/scripts/htmx/extensions.ts"),
		scripts_template.HtmxExtensionsTSTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/scripts/htmx/index.ts"),
		scripts_template.HtmxIndexTSTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	//! /resources/css

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/css/main.css"),
		css_template.MainCssTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/resources/css/tailwind.css"),
		css_template.TailwindCssTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.GoImports(projectPath)
	if err != nil {
		return err
	}

	err = utils.GoFumpt(projectPath)
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

	err = utils.PnpmInstall(projectPath)
	if err != nil {
		return err
	}

	return nil
}
