package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(linkControllerCmd)
}

var linkControllerCmd = &cobra.Command{
	Use:   "link:controller (ControllerName)",
	Short: "Add in main app container a controller",
	Args:  cobra.MinimumNArgs(0),
	Run:   LinkController,
}

func LinkController(cmd *cobra.Command, args []string) {
	controllers := analyzer.AnalyzeProjectControllers()

	choices := []form.Choice[string]{}

	for _, controller := range controllers {
		choices = append(choices, form.Choice[string]{
			Name:  controller.Name,
			Value: controller.Name,
		})
	}

	var controllerName string
	if len(args) == 0 {
		controllerName = form.AskChoiceSearch("Controller to link :", choices)
	} else {
		controllerName = args[0]
	}

	var selectedController *analyzer.AnalyzedController

	for _, controller := range controllers {
		if controllerName != controller.Name {
			continue
		}

		selectedController = &controller
		break
	}

	if selectedController == nil {
		logger.MainLogger.Fatalf("Can't find the controller \"%s\"", controllerName)
	}

	err := actions.LinkController(controllerName, *selectedController)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Link Controller : %v", err)
	}
}
