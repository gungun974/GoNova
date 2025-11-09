package command

import (
	"fmt"
	"os"

	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().String("database", "none", "Choose the database module: sqlite | postgres | none")
	initCmd.Flags().Bool("no-git", false, "Skip Git initialization for the project")
	initCmd.Flags().Bool("nix", false, "Include the Nix module")
}

var initCmd = &cobra.Command{
	Use:   "init (name)",
	Short: "Init Nova Core in current directory",
	Args:  cobra.MinimumNArgs(0),
	Run:   InitNova,
}

func InitNova(cmd *cobra.Command, args []string) {
	isInteractive := term.IsTerminal(int(os.Stdout.Fd()))

	databaseFlag, _ := cmd.Flags().GetString("database")
	noGitFlag, _ := cmd.Flags().GetBool("no-git")
	nixFlag, _ := cmd.Flags().GetBool("nix")

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

	projectName := ""
	if len(args) == 0 {
		projectName = form.AskInputWithPlaceholder(
			"The golang module name :",
			"github.com/user/project",
		)
	} else {
		projectName = args[0]
	}

	var database string

	if !dbFlagSet {
		database = form.AskChoice("Do you want to install a Database module :", []form.Choice[string]{
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
				Value: "postgres",
			},
		})
	} else {
		if databaseFlag != "sqlite" && databaseFlag != "postgres" && databaseFlag != "none" {
			fmt.Printf("Error: invalid value for --database: %q\n", databaseFlag)
			fmt.Println("Allowed values: sqlite | postgres | none")
			fmt.Println()
			_ = cmd.Help()
			os.Exit(1)
		}
		database = databaseFlag
	}

	var enableNoGit bool

	if !noGitFlagSet {
		enableNoGit = !form.AskOption("Do you want to initialize Git :", true, "Yes", "No")
	} else {
		enableNoGit = noGitFlag
	}

	var enableNix bool

	if !nixFlagSet {
		enableNix = form.AskOption("Do you want to install the Nix module :", true, "Yes", "No")
	} else {
		enableNix = nixFlag
	}

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

	if database == "postgres" {
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
