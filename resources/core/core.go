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

//go:embed .air.toml.tmpl
var airFile []byte

var AirTemplate = template.Must(template.New(".air.toml").Parse(string(airFile)))

//go:embed .env.tmpl
var envFile []byte

var EnvTemplate = template.Must(template.New(".env").Parse(string(envFile)))
