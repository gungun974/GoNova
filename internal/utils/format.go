package utils

import (
	"os"

	"github.com/gungun974/gonova/internal/logger"
	"golang.org/x/tools/imports"
)

func GoImportFormatFile(path string) error {
	formatted, err := imports.Process(path, nil, nil)
	if err != nil {
		logger.CommandLogger.Errorf("Failed to use goimports in file %s : %v", path, err)
		return err
	}

	err = os.WriteFile(path, []byte(formatted), 0)
	if err != nil {
		logger.CommandLogger.Errorf("Failed to write in file %s : %v", path, err)
		return err
	}

	return nil
}
