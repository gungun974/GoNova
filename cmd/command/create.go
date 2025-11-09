package command

import (
	"fmt"
	"os"
	"path"

	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().String("database", "none", "Choose the database module: sqlite | postgres | none")
	createCmd.Flags().Bool("no-git", false, "Skip Git initialization for the project")
	createCmd.Flags().Bool("nix", false, "Include the Nix module")
}

var createCmd = &cobra.Command{
	Use:   "create [directory-name] (name)",
	Short: "Create Nova Core in a specific directory",
	Args:  cobra.MinimumNArgs(1),
	Run:   CreateNova,
}

func CreateNova(cmd *cobra.Command, args []string) {
	isInteractive := term.IsTerminal(int(os.Stdout.Fd()))

	dbFlagSet := cmd.Flags().Changed("database")
	noGitFlagSet := cmd.Flags().Changed("no-git")
	nixFlagSet := cmd.Flags().Changed("nix")

	if !isInteractive {
		if !dbFlagSet || !noGitFlagSet || !nixFlagSet {
			fmt.Println("Error: In non-interactive mode, all flags must be provided explicitly.")
			if len(args) == 0 {
				fmt.Println("Error: In non-interactive mode, all arguments must be provided explicitly.")
			}
			fmt.Println()
			_ = cmd.Help()
			os.Exit(1)
		}
		if len(args) == 0 {
			fmt.Println("Error: In non-interactive mode, all arguments must be provided explicitly.")
			fmt.Println()
			_ = cmd.Help()
			os.Exit(1)
		}
	}

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
