package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	installCmd := &cobra.Command{Use: "install"}

	installCmd.AddCommand(installPostgreDatabaseCmd)
	installCmd.AddCommand(installSqliteDatabaseCmd)

	rootCmd.AddCommand(installCmd)
}

var installPostgreDatabaseCmd = &cobra.Command{
	Use:   "postgre",
	Short: "Install Nova Postgree Database module in current directory",
	Run:   InstallPostgreDatabase,
}

func InstallPostgreDatabase(_ *cobra.Command, _ []string) {
	err := actions.InstallPostgreDatabase()
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Install Postgre Database : %v", err)
	}
}

var installSqliteDatabaseCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "Install Nova Sqlite Database module in current directory",
	Run:   InstallSqliteDatabase,
}

func InstallSqliteDatabase(_ *cobra.Command, _ []string) {
	err := actions.InstallSqliteDatabase()
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Install Sqlite Database : %v", err)
	}
}
