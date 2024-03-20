package utils

import (
	"os"
	"os/exec"
	"strings"

	"github.com/gungun974/gonova/internal/logger"
)

func VerifyCmd(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func ExecuteCmd(name string, args []string, dir string) error {
	command := exec.Command(name, args...)
	command.Dir = dir

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	logger.CommandLogger.Infof("Run : %s %s", name, strings.Join(args, " "))

	if err := command.Run(); err != nil {
		logger.CommandLogger.Errorf("Fail : %s %s", name, strings.Join(args, " "))
		return err
	}

	return nil
}
