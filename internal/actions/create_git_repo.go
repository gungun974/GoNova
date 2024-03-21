package actions

import (
	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
)

func CheckCanCreateGitRepo() error {
	nameSet, err := utils.CheckGitConfig("user.name")
	if err != nil {
		return err
	}
	if !nameSet {
		logger.MainLogger.Error("user.name is not set in git config.")
		logger.MainLogger.Fatal("Please set up git config before trying again.")
	}

	emailSet, err := utils.CheckGitConfig("user.email")
	if err != nil {
		return err
	}
	if !emailSet {
		logger.MainLogger.Error("user.email is not set in git config.")
		logger.MainLogger.Fatal("Please set up git config before trying again.")
	}

	return nil
}

func CreateGitRepo() error {
	projectPath := "."

	logger.MainLogger.Info("Create Git Repo")

	err := CheckCanCreateGitRepo()
	if err != nil {
		return err
	}

	err = utils.InitGit(projectPath)
	if err != nil {
		return err
	}

	err = utils.GitAddAllFiles(projectPath)
	if err != nil {
		return err
	}

	err = utils.CreateGitInitialCommit(projectPath)
	if err != nil {
		return err
	}

	return nil
}
