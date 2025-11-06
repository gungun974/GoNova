package core_template

import (
	_ "embed"
	"text/template"
)

//go:embed Makefile.tmpl
var makefileFile []byte

var MakefileTemplate = template.Must(template.New("Makefile").Parse(string(makefileFile)))

//go:embed .golangci.yml.tmpl
var golangciFile []byte

var GolangciTemplate = template.Must(template.New(".golangci.yml").Parse(string(golangciFile)))

//go:embed .gitignore.tmpl
var gitignoreFile []byte

var GitignoreTemplate = template.Must(template.New(".gitignore").Parse(string(gitignoreFile)))

//go:embed .env.tmpl
var envFile []byte

var EnvTemplate = template.Must(template.New(".env").Parse(string(envFile)))

//! Vite / JS / CSS

//go:embed eslint.config.mjs.tmpl
var eslintConfigFile []byte

var EslintConfigFileTemplate = template.Must(
	template.New("eslint.config.mjs").Parse(string(eslintConfigFile)),
)

//go:embed .npmrc.tmpl
var npmrcFile []byte

var NpmrcTemplate = template.Must(template.New(".npmrc").Parse(string(npmrcFile)))

//go:embed package.json.tmpl
var packageJsonFile []byte

var PackageJsonTemplate = template.Must(template.New("package.json").Parse(string(packageJsonFile)))

//go:embed pnpm-workspace.yaml.tmpl
var pnpmWorkspaceFile []byte

var PnpmWorkspaceTemplate = template.Must(
	template.New("pnpm-workspace.yaml").Parse(string(pnpmWorkspaceFile)),
)

//go:embed prettier.config.js.tmpl
var prettierFile []byte

var PrettierTemplate = template.Must(template.New("prettier.config.js").Parse(string(prettierFile)))

//go:embed tsconfig.json.tmpl
var tsconfigFile []byte

var TsconfigTemplate = template.Must(template.New("tsconfig.json").Parse(string(tsconfigFile)))

//go:embed vite.config.ts.tmpl
var viteFile []byte

var ViteTemplate = template.Must(template.New("vite.config.ts").Parse(string(viteFile)))
