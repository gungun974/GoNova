package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("postgre", "", false, "Init with postgre module")
	initCmd.Flags().BoolP("sqlite", "", false, "Init with sqlite sqlite")
}

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Init Nova Core in current directory",
	Args:  cobra.MinimumNArgs(1),
	Run:   InitNova,
}

func InitNova(cmd *cobra.Command, args []string) {
	enablePostgre, _ := cmd.Flags().GetBool("postgre")
	enableSqlite, _ := cmd.Flags().GetBool("sqlite")

	if enablePostgre && enableSqlite {
		logger.MainLogger.Fatal("You can't install postgree and sqlite at both time")
	}

	err := actions.InstallCore(args[0], enablePostgre, enableSqlite)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Install Core : %v", err)
	}

	if enablePostgre {
		err := actions.InstallPostgreDatabase()
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Install Postgre Database : %v", err)
		}
	}

	if enableSqlite {
		err := actions.InstallSqliteDatabase()
		if err != nil {
			logger.MainLogger.Fatalf("Failed to Install Sqlite Database : %v", err)
		}
	}

	err = actions.FinishInstall()
	if err != nil {
		logger.MainLogger.Fatalf("Failed to finish Install : %v", err)
	}
}
