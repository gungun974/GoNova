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
		routerName = form.AskInputWithPlaceholder("The router name :", "Post")
	} else {
		routerName = args[0]
	}

	urlMountPath := form.AskInputWithPlaceholder("Url mount path :", "/post")

	err := actions.MakeRouter(routerName, urlMountPath)
	if err != nil {
		logger.MainLogger.Fatalf("Failed to Make Router : %v", err)
	}
}
