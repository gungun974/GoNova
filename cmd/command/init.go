package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Init Nova Core in current directory",
	Args:  cobra.MinimumNArgs(1),
	Run:   InitNova,
}

func InitNova(_ *cobra.Command, args []string) {
	err := actions.InstallCore(args[0])
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Install Core : %v", err)
	}
}
