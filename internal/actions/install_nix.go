package actions

import (
	"os"
	"path/filepath"

	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	nix_template "github.com/gungun974/gonova/resources/nix"
)

func InstallNix() error {
	projectPath := "."

	projectName, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	if _, err := os.Stat(filepath.Join(projectPath, "/flake.nix")); err == nil {
		logger.MainLogger.Fatal("Can't Install Nix Module when flake.nix is already installed")
	}

	if !utils.VerifyCmd("go") {
		logger.MainLogger.Fatal("Can't Install Nix without go in PATH")
	}

	tmplVersion, err := utils.GetTemplVersion(projectPath)
	if err != nil {
		logger.MainLogger.Fatal(err)
	}

	logger.MainLogger.Info("Install Nix")

	projectGlobalTemplateConfig := struct {
		ProjectName string
		TmplVersion string
	}{
		ProjectName: projectName,
		TmplVersion: tmplVersion,
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/flake.nix"),
		nix_template.FlakeNixTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/nix/helpers/esbuild/default.nix"),
		nix_template.EsbuildNixTemplate,
		projectGlobalTemplateConfig,
	)
	if err != nil {
		return err
	}

	err = utils.CreateFileFromTemplate(
		filepath.Join(projectPath, "/.envrc"),
		nix_template.EnvrcTemplate,
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

	err = utils.NixFlakeUpdate(projectPath)
	if err != nil {
		logger.MainLogger.Warn("Can't update nix flake")
	}

	return nil
}
