package utils

func NixFlakeUpdate(appDir string) error {
	if err := ExecuteCmd("nix",
		[]string{"flake", "update"},
		appDir); err != nil {
		return err
	}

	return nil
}
