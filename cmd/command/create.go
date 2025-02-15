package command

import (
	"os"
	"path"

	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create [directory-name] (name)",
	Short: "Create Nova Core in a specific directory",
	Args:  cobra.MinimumNArgs(1),
	Run:   CreateNova,
}

func CreateNova(cmd *cobra.Command, args []string) {
	err := utils.CreateDirectory(args[0])
	if err != nil {
		logger.MainLogger.Fatalf("Failed to create the project directory : %v", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		logger.MainLogger.Fatalf("Failed to get current directory : %v", err)
	}

	err = os.Chdir(path.Join(currentDir, args[0]))
	if err != nil {
		logger.MainLogger.Fatalf("Failed to switch to the project directory : %v", err)
	}

	InitNova(cmd, args[1:])

	logger.MainLogger.Info("✨ Nova project has been created.")
}
