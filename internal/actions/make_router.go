package actions

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/injector"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	make_router_template "github.com/gungun974/gonova/resources/make/router"
)

func MakeRouter(routerName string, urlMountPath string) error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	newRouteFilePath := fmt.Sprintf("/internal/routes/%s.go", helpers.ToSnakeCase(routerName))

	if _, err := os.Stat(filepath.Join(projectPath, newRouteFilePath)); err == nil {
		logger.MainLogger.Fatal(
			"Can't Make Router when a route with the same name is already present",
		)
	}

	logger.MainLogger.Info("Make Router")

	projectGlobalTemplateConfig := struct {
		ProjectName string
		RouterName  string
		RouteURL    string
	}{
		ProjectName: projectName,
		RouterName:  helpers.CapitalizeFirstLetter(routerName) + "Router",
		RouteURL:    path.Clean("/" + urlMountPath),
	}

	//! /internal/routes

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, newRouteFilePath),
		make_router_template.RouteGoTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	injector.InjectRoutesNewRouter(
		filepath.Join(projectPath, "/internal/routes/routes.go"),
		projectGlobalTemplateConfig.RouterName,
		projectGlobalTemplateConfig.RouteURL,
	)

	err = utils.GoFmt(projectPath)
	if err != nil {
		return err
	}

	return nil
}
