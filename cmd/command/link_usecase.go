package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(linkUsecaseCmd)
}

var linkUsecaseCmd = &cobra.Command{
	Use:   "link:usecase (UsecaseName)",
	Short: "Link a usecase with an existing controller",
	Args:  cobra.MinimumNArgs(0),
	Run:   LinkUsecase,
}

func LinkUsecase(cmd *cobra.Command, args []string) {
	usecases := analyzer.AnalyzeProjectUsecases()

	usecaseChoices := []form.Choice[string]{}

	for _, usecase := range usecases {
		usecaseChoices = append(usecaseChoices, form.Choice[string]{
			Name:  usecase.Name,
			Value: usecase.Name,
		})
	}

	var usecaseName string
	if len(args) == 0 {
		usecaseName = form.AskChoiceSearch("Usecase to link :", usecaseChoices)
	} else {
		usecaseName = args[0]
	}

	var selectedUsecase *analyzer.AnalyzedUsecase

	for _, usecase := range usecases {
		if usecaseName != usecase.Name {
			continue
		}

		selectedUsecase = &usecase
		break
	}

	if selectedUsecase == nil {
		logger.MainLogger.Fatalf("Can't find the usecase \"%s\"", usecaseName)
	}

	controllers := analyzer.AnalyzeProjectControllers()

	controllerChoices := []form.Choice[string]{}

	for _, controller := range controllers {
		controllerChoices = append(controllerChoices, form.Choice[string]{
			Name:  controller.Name,
			Value: controller.Name,
		})
	}

	var controllerName string
	if len(args) <= 1 {
		controllerName = form.AskChoiceSearch("Controller to link :", controllerChoices)
	} else {
		controllerName = args[1]
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

	err := actions.LinkUsecase(*selectedUsecase, *selectedController)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Link Usecase : %v", err)
	}
}
