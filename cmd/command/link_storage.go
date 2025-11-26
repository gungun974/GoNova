package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(linkStorageCmd)
}

var linkStorageCmd = &cobra.Command{
	Use:   "link:storage (StorageName) (UsecaseName)",
	Short: "Link a storage with an existing usecase",
	Args:  cobra.MinimumNArgs(0),
	Run:   LinkStorage,
}

func LinkStorage(cmd *cobra.Command, args []string) {
	storages := analyzer.AnalyzeProjectStorages()

	storageChoices := []form.Choice[string]{}

	for _, storage := range storages {
		storageChoices = append(storageChoices, form.Choice[string]{
			Name:  storage.Name,
			Value: storage.Name,
		})
	}

	var storageName string
	if len(args) == 0 {
		storageName = form.AskChoiceSearch("Storage to link :", storageChoices)
	} else {
		storageName = args[0]
	}

	var selectedStorage *analyzer.AnalyzedStorage

	for _, storage := range storages {
		if storageName != storage.Name {
			continue
		}

		selectedStorage = &storage
		break
	}

	if selectedStorage == nil {
		logger.MainLogger.Fatalf("Can't find the storage \"%s\"", storageName)
	}

	repositories := analyzer.AnalyzeProjectRepositories()
	presenters := analyzer.AnalyzeProjectPresenters()
	usecases := analyzer.AnalyzeProjectUsecases(repositories, storages, presenters)

	usecaseChoices := []form.Choice[string]{}

usecase_loop:
	for _, usecase := range usecases {
		for _, dependency := range usecase.Dependencies {
			if dependency == nil {
				continue
			}
			if dependency.GetName() == selectedStorage.GetName() {
				continue usecase_loop
			}
		}
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

usecase_loop2:
	for _, usecase := range usecases {
		for _, dependency := range usecase.Dependencies {
			if dependency == nil {
				continue
			}
			if dependency.GetName() == selectedStorage.GetName() {
				continue usecase_loop2
			}
		}
		if usecaseName != usecase.Name {
			continue
		}

		selectedUsecase = &usecase
		break
	}

	if selectedUsecase == nil {
		logger.MainLogger.Fatalf("Can't find the usecase \"%s\"", usecaseName)
	}

	err := actions.LinkStorage(*selectedStorage, *selectedUsecase)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Link Storage : %v", err)
	}
}
