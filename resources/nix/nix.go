package nix_template

import (
	_ "embed"
	"text/template"
)

//go:embed flake.nix.tmpl
var flakeNixFile []byte

var FlakeNixTemplate = template.Must(template.New("flake.nix").Parse(string(flakeNixFile)))

//go:embed helpers/esbuild/default.nix.tmpl
var esbuildNixFile []byte

var EsbuildNixTemplate = template.Must(
	template.New("helpers/esbuild/default.nix").Parse(string(esbuildNixFile)),
)

//go:embed .envrc.tmpl
var envrcFile []byte

var EnvrcTemplate = template.Must(template.New(".envrc").Parse(string(envrcFile)))
