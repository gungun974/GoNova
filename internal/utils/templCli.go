package utils

func TemplGenerate(appDir string) error {
	if err := ExecuteCmd("go",
		[]string{"run", "github.com/a-h/templ/cmd/templ", "generate", "-v"},
		appDir); err != nil {
		return err
	}

	return nil
}
