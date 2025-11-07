package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(makeRepositoryCmd)
}

var makeRepositoryCmd = &cobra.Command{
	Use:   "make:repository (RepositoryName)",
	Short: "Create a new Repository",
	Args:  cobra.MinimumNArgs(0),
	Run:   MakeRepository,
}

func MakeRepository(cmd *cobra.Command, args []string) {
	repositoryName := ""
	if len(args) == 0 {
		repositoryName = form.AskInputWithPlaceholder("The Repository name :", "Post")
	} else {
		repositoryName = args[0]
	}

	err := actions.MakeRepository(repositoryName)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Make Repository : %v", err)
	}
}
