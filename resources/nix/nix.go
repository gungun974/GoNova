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

var EsbuildNixTemplate = template.Must(template.New("helpers/esbuild/default.nix").Parse(string(esbuildNixFile)))

//go:embed helpers/fetch-pnpm-deps/default.nix.tmpl
var pnpmNixFile []byte

var PnpmNixTemplate = template.Must(template.New("helpers/fetch-pnpm-deps/default.nix").Parse(string(pnpmNixFile)))

//go:embed helpers/fetch-pnpm-deps/pnpm-config-hook.sh.tmpl
var pnpmConfigHookFile []byte

var PnpmConfigHookTemplate = template.Must(template.New("helpers/fetch-pnpm-deps/pnpm-config-hook.sh").Parse(string(pnpmConfigHookFile)))

//go:embed .envrc.tmpl
var envrcFile []byte

var EnvrcTemplate = template.Must(template.New(".envrc").Parse(string(envrcFile)))
