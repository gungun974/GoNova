package command

import (
	"cmp"
	"slices"
	"strings"

	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(makeRepositoryCmd)
}

var makeRepositoryCmd = &cobra.Command{
	Use:   "make:repository (RepositoryName) (ModelName)",
	Short: "Create a new Repository with a model or merge a model in a existing repository",
	Args:  cobra.MinimumNArgs(0),
	Run:   MakeRepository,
}

func MakeRepository(cmd *cobra.Command, args []string) {
	entities := analyzer.AnalyzeProjectEntities()
	models := analyzer.AnalyzeProjectModels(entities)
	repositories := analyzer.AnalyzeProjectRepositories()

	modelChoices := []form.Choice[*string]{
		{
			Name:  "No model",
			Value: nil,
		},
	}

model_loop:
	for _, model := range models {
		for _, entity := range entities {
			if entity.Equal(model.Entity) {
				modelChoices = append(modelChoices, form.Choice[*string]{
					Name:  model.Entity.Name,
					Value: &model.Name,
				})
				continue model_loop
			}
		}
	}

	var modelName *string
	if len(args) <= 1 {
		modelName = form.AskChoiceSearchWithAlwaysShowFirstChoice("Based repository on model :", modelChoices)
	} else {
		modelName = &args[1]
	}

	var selectedModel *analyzer.AnalyzedModel

	if modelName != nil {
		for _, model := range models {
			if *modelName != model.Name {
				continue
			}

			selectedModel = &model

			break
		}
	}

	if modelName != nil && selectedModel == nil {
		logger.MainLogger.Fatalf("Can't find the model \"%s\"", *modelName)
		return
	}

	var repositoryName *string
	if len(args) == 0 {
		if selectedModel != nil {
			repositoryChoices := []form.Choice[*string]{
				{
					Name:  "New repository",
					Value: nil,
				},
			}

			slices.SortFunc(repositories, func(a, b analyzer.AnalyzedRepository) int {
				aHasPrefix := strings.HasPrefix(a.Name, selectedModel.Entity.Name)
				bHasPrefix := strings.HasPrefix(b.Name, selectedModel.Entity.Name)

				if aHasPrefix && !bHasPrefix {
					return -1
				}
				if !aHasPrefix && bHasPrefix {
					return 1
				}

				return cmp.Compare(a.Name, b.Name)
			})

			for _, repository := range repositories {
				baseName := strings.TrimSuffix(repository.Name, "Repository")
				repositoryChoices = append(repositoryChoices, form.Choice[*string]{
					Name:  repository.Name,
					Value: &baseName,
				})
			}

			repositoryName = form.AskChoiceSearchWithAlwaysShowFirstChoice("Repository :", repositoryChoices)

			if repositoryName == nil {
				name := form.AskInputWithPlaceholder("Repository name :", "Post")
				repositoryName = &name
			}
		} else {
			name := form.AskInputWithPlaceholder("Repository name :", "Post")
			repositoryName = &name

		}
	} else {
		repositoryName = &args[0]
	}

	err := actions.MakeRepository(*repositoryName, selectedModel)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Make Repository : %v", err)
	}
}
