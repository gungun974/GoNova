package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(linkPresenterCmd)
}

var linkPresenterCmd = &cobra.Command{
	Use:   "link:presenter (PresenterName) (UsecaseName)",
	Short: "Link a presenter with an existing usecase",
	Args:  cobra.MinimumNArgs(0),
	Run:   LinkPresenter,
}

func LinkPresenter(cmd *cobra.Command, args []string) {
	presenters := analyzer.AnalyzeProjectPresenters()

	presenterChoices := []form.Choice[string]{}

	for _, presenter := range presenters {
		presenterChoices = append(presenterChoices, form.Choice[string]{
			Name:  presenter.Name,
			Value: presenter.Name,
		})
	}

	var presenterName string
	if len(args) <= 1 {
		presenterName = form.AskChoiceSearch("Presenter to link :", presenterChoices)
	} else {
		presenterName = args[1]
	}

	var selectedPresenter *analyzer.AnalyzedPresenter

	for _, presenter := range presenters {
		if presenterName != presenter.Name {
			continue
		}

		selectedPresenter = &presenter
		break
	}

	if selectedPresenter == nil {
		logger.MainLogger.Fatalf("Can't find the presenter \"%s\"", presenterName)
	}

	repositories := analyzer.AnalyzeProjectRepositories()
	usecases := analyzer.AnalyzeProjectUsecases(repositories, presenters)

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

	err := actions.LinkPresenter(*selectedPresenter, *selectedUsecase)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Link Presenter : %v", err)
	}
}
