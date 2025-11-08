package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(makePresenterCmd)
}

var makePresenterCmd = &cobra.Command{
	Use:   "make:presenter (PresenterName)",
	Short: "Create a new Presenter",
	Args:  cobra.MinimumNArgs(0),
	Run:   MakePresenter,
}

func MakePresenter(cmd *cobra.Command, args []string) {
	presenterName := ""
	if len(args) == 0 {
		presenterName = form.AskInputWithPlaceholder("The Presenter name :", "Post")
	} else {
		presenterName = args[0]
	}

	newPresenter, err := actions.MakePresenter(presenterName)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Make Presenter : %v", err)
	}

	if form.AskOption("Do you want to link with a usecase :", true, "Link", "No") {
		LinkPresenter(cmd, []string{
			newPresenter,
		})
	}
}
