package command

import (
	"github.com/gungun974/gonova/internal/actions"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(makeRouterCmd)
}

var makeRouterCmd = &cobra.Command{
	Use:   "make:router (RouterName)",
	Short: "Create a new router with sub mount",
	Args:  cobra.MinimumNArgs(0),
	Run:   MakeRouter,
}

func MakeRouter(cmd *cobra.Command, args []string) {
	routerName := ""
	if len(args) == 0 {
		routerName = form.AskInput("The router name :")
	} else {
		routerName = args[0]
	}

	urlMountPath := form.AskInput("Url mount path :")

	err := actions.MakeRouter(routerName, urlMountPath)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Make Router : %v", err)
	}
}
