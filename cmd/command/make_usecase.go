package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(makeUsecaseCmd)
}

var makeUsecaseCmd = &cobra.Command{
	Use:   "make:usecase (UsecaseName)",
	Short: "Create a new Usecase",
	Args:  cobra.MinimumNArgs(0),
	Run:   MakeUsecase,
}

func MakeUsecase(cmd *cobra.Command, args []string) {
	usecaseName := ""
	if len(args) == 0 {
		usecaseName = form.AskInputWithPlaceholder("The Usecase name :", "Post")
	} else {
		usecaseName = args[0]
	}

	err := actions.MakeUsecase(usecaseName)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Make Usecase : %v", err)
	}
}
