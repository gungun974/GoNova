package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(makeEntityCmd)
}

var makeEntityCmd = &cobra.Command{
	Use:   "make:entity (EntityName)",
	Short: "Create a new Entity",
	Args:  cobra.MinimumNArgs(0),
	Run:   MakeEntity,
}

func MakeEntity(cmd *cobra.Command, args []string) {
	entityName := ""
	if len(args) == 0 {
		entityName = form.AskInputWithPlaceholder("The Entity name :", "Post")
	} else {
		entityName = args[0]
	}

	err := actions.MakeEntity(entityName)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Make Entity : %v", err)
	}
}
