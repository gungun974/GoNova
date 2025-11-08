package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(linkRepositoryCmd)
}

var linkRepositoryCmd = &cobra.Command{
	Use:   "link:repository (RepositoryName) (UsecaseName)",
	Short: "Link a repository with an existing usecase",
	Args:  cobra.MinimumNArgs(0),
	Run:   LinkRepository,
}

func LinkRepository(cmd *cobra.Command, args []string) {
	repositories := analyzer.AnalyzeProjectRepositories()

	repositoryChoices := []form.Choice[string]{}

	for _, repository := range repositories {
		repositoryChoices = append(repositoryChoices, form.Choice[string]{
			Name:  repository.Name,
			Value: repository.Name,
		})
	}

	var repositoryName string
	if len(args) == 0 {
		repositoryName = form.AskChoiceSearch("Repository to link :", repositoryChoices)
	} else {
		repositoryName = args[0]
	}

	var selectedRepository *analyzer.AnalyzedRepository

	for _, repository := range repositories {
		if repositoryName != repository.Name {
			continue
		}

		selectedRepository = &repository
		break
	}

	if selectedRepository == nil {
		logger.MainLogger.Fatalf("Can't find the repository \"%s\"", repositoryName)
	}

	presenters := analyzer.AnalyzeProjectPresenters()
	usecases := analyzer.AnalyzeProjectUsecases(repositories, presenters)

	usecaseChoices := []form.Choice[string]{}

	for _, usecase := range usecases {
		usecaseChoices = append(usecaseChoices, form.Choice[string]{
			Name:  usecase.Name,
			Value: usecase.Name,
		})
	}

	var usecaseName string
	if len(args) <= 1 {
		usecaseName = form.AskChoiceSearch("Usecase to link :", usecaseChoices)
	} else {
		usecaseName = args[1]
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

	err := actions.LinkRepository(*selectedRepository, *selectedUsecase)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Link Repository : %v", err)
	}
}
