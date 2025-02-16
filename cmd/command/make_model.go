package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/analyzer"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(makeModelCmd)
}

var makeModelCmd = &cobra.Command{
	Use:   "make:model",
	Short: "Create a new Model",
	Args:  cobra.MinimumNArgs(0),
	Run:   MakeModel,
}

func MakeModel(cmd *cobra.Command, args []string) {
	entities := analyzer.AnalyzeProjectEntities()

	choices := []form.Choice{}

	for _, entity := range entities {
		choices = append(choices, form.Choice{
			Name:  entity.Name,
			Value: entity.Name,
		})
	}

	entityName := ""
	if len(args) == 0 {
		entityName = form.AskChoiceSearch("Based model on entity :", choices)
	} else {
		entityName = args[0]
	}

	var selectedEntity *analyzer.AnalyzedEntity

	for _, entity := range entities {
		if entityName != entity.Name {
			continue
		}

		selectedEntity = &entity

		break
	}

	if selectedEntity == nil {
		logger.MainLogger.Logger.Fatalf("Can't find the entity \"%s\"", entityName)
		return
	}

	err := actions.MakeModel(*selectedEntity)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Make Model : %v", err)
	}
}
