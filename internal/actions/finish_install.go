package actions

import (
	"github.com/gungun974/gonova/internal/utils"
)

func FinishInstall() error {
	projectPath := "."

	err := utils.GoTidy(projectPath)
	if err != nil {
		return err
	}

	return nil
}
