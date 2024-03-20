package utils

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/gungun974/gonova/internal/logger"
)

func CreateDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0751)
		if err != nil {
			logger.MainLogger.Errorf("Failed to create directory : %v", path)
			return err
		}
	}

	return nil
}

func CreateFileFromTemplate(path string, template *template.Template, data any) error {
	directory := filepath.Dir(path)

	if err := CreateDirectory(directory); err != nil {
		logger.MainLogger.Errorf("Failed to create directory for file : %v", path)
		return err
	}

	createdFile, err := os.Create(path)
	if err != nil {
		logger.MainLogger.Errorf("Failed to create file : %v", path)
		return err
	}

	defer createdFile.Close()

	return template.Execute(createdFile, data)
}
