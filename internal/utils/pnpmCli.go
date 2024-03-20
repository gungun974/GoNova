package utils

func PnpmInstall(appDir string) error {
	if err := ExecuteCmd("pnpm",
		[]string{"install"},
		appDir); err != nil {
		return err
	}

	return nil
}
