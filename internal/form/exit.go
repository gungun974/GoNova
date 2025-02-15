package form

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gungun974/gonova/internal/logger"
)

func handleExit(tprogram *tea.Program, hasExit bool) {
	if hasExit {
		if err := tprogram.ReleaseTerminal(); err != nil {
			logger.MainLogger.Fatal(err)
		}
		os.Exit(1)
	}
}
