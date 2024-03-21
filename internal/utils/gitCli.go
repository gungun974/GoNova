package utils

import "os/exec"

func CheckGitConfig(key string) (bool, error) {
	cmd := exec.Command("git", "config", "--get", key)
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// The command failed to run.
			if exitError.ExitCode() == 1 {
				// The 'git config --get' command returns 1 if the key was not found.
				return false, nil
			}
		}
		// Some other error occurred.
		return false, err
	}
	// The command ran successfully, so the key is set.
	return true, nil
}

func InitGit(appDir string) error {
	if err := ExecuteCmd("git",
		[]string{"init"},
		appDir); err != nil {
		return err
	}

	return nil
}

func GitAddAllFiles(appDir string) error {
	if err := ExecuteCmd("git",
		[]string{"add", "."},
		appDir); err != nil {
		return err
	}

	return nil
}

func CreateGitInitialCommit(appDir string) error {
	if err := ExecuteCmd("git",
		[]string{"commit", "-m", "Initial commit"},
		appDir); err != nil {
		return err
	}

	return nil
}
