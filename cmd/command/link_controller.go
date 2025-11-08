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
	repositories := analyzer.AnalyzeProjectRepositories()
	presenters := analyzer.AnalyzeProjectPresenters()
	usecases := analyzer.AnalyzeProjectUsecases(repositories, presenters)
	controllers := analyzer.AnalyzeProjectControllers(usecases)
	container := analyzer.AnalyzeProjectContainer(controllers)

	choices := []form.Choice[string]{}

controller_loop:
	for _, controller := range controllers {
		for _, dependency := range container.Dependencies {
			if dependency.GetName() == controller.GetName() {
				continue controller_loop
			}
		}
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

controller_loop2:
	for _, controller := range controllers {
		for _, dependency := range container.Dependencies {
			if dependency.GetName() == controller.GetName() {
				continue controller_loop2
			}
		}
		if controllerName != controller.Name {
			continue
		}

		selectedController = &controller
		break
	}

	if selectedController == nil {
		logger.MainLogger.Fatalf("Can't find the controller \"%s\"", controllerName)
	}

	err := actions.LinkController(*selectedController)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Link Controller : %v", err)
	}
}
