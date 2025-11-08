package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(makeControllerCmd)
}

var makeControllerCmd = &cobra.Command{
	Use:   "make:controller (ControllerName)",
	Short: "Create a new Controller",
	Args:  cobra.MinimumNArgs(0),
	Run:   MakeController,
}

func MakeController(cmd *cobra.Command, args []string) {
	controllerName := ""
	if len(args) == 0 {
		controllerName = form.AskInputWithPlaceholder("Controller name :", "Post")
	} else {
		controllerName = args[0]
	}

	newController, err := actions.MakeController(controllerName)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Make Controller : %v", err)
	}

	if form.AskOption("Do you want to link with the container :", true, "Link", "No") {
		LinkController(cmd, []string{
			newController,
		})
	}
}
