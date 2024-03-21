package actions

import (
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
)

func FinishInstall() error {
	projectPath := "."

	logger.MainLogger.Info("Create Git Repo")

	err := utils.GoTidy(projectPath)
	if err != nil {
		return err
	}

	return nil
}
