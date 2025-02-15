package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init (name)",
	Short: "Init Nova Core in current directory",
	Args:  cobra.MinimumNArgs(0),
	Run:   InitNova,
}

func InitNova(cmd *cobra.Command, args []string) {
	projectName := ""
	if len(args) == 0 {
		projectName = form.AskInputWithPlaceholder(
			"The golang module name :",
			"github.com/user/project",
		)
	} else {
		projectName = args[0]
	}

	database := form.AskChoice("Do you want to install a Database module :", []form.Choice{
		{
			Name:  "No",
			Value: "none",
		},
		{
			Name:  "SQLite",
			Value: "sqlite",
		},
		{
			Name:  "PostgreSQL",
			Value: "postgre",
		},
	})

	enableNoGit := !form.AskOption("Do you want to initialize Git :", true, "Yes", "No")

	enableNix := form.AskOption("Do you want to install the Nix module :", true, "Yes", "No")

	if !enableNoGit {
		err := actions.CheckCanCreateGitRepo()
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Check Git : %v", err)
		}

	}

	err := actions.InstallCore(projectName)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Install Core : %v", err)
	}

	if database == "postgre" {
		err := actions.InstallPostgreDatabase()
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Install Postgre Database : %v", err)
		}
	}

	if database == "sqlite" {
		err := actions.InstallSqliteDatabase()
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Install Sqlite Database : %v", err)
		}
	}

	if enableNix {
		err := actions.InstallNix()
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Install Sqlite Database : %v", err)
		}
	}

	err = actions.FinishInstall()
	if err != nil {
		logger.MainLogger.Fatalf("Failed to finish Install : %v", err)
	}

	if !enableNoGit {
		err := actions.CreateGitRepo()
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Create Git Repo : %v", err)
		}

	}
}
