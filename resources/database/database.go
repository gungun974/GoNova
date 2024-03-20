package database_template

import (
	_ "embed"
	"text/template"
)

//go:embed EMPTY_MIGRATION.sql.tmpl
var emptyMigrationFile []byte

var EmptyMigrationTemplate = template.Must(template.New("EMPTY_MIGRATION.sql").Parse(string(emptyMigrationFile)))

//go:embed logger.go.tmpl
var loggerGoFile []byte

var LoggerGoTemplate = template.Must(template.New("logger.go").Parse(string(loggerGoFile)))

//go:embed postgre/database.go.tmpl
var postgreDatabaseGoFile []byte

var PostgreDatabaseGoTemplate = template.Must(template.New("postgre/database.go").Parse(string(postgreDatabaseGoFile)))

//go:embed postgre/.env.tmpl
var postgreEnvFile []byte

var PostgreEnvFileTemplate = template.Must(template.New("postgre/.env").Parse(string(postgreEnvFile)))

//go:embed postgre/docker-compose.yml.tmpl
var postgreDockerComposeFile []byte

var PostgreDockerComposeFileTemplate = template.Must(template.New("postgre/docker-compose.yml").Parse(string(postgreDockerComposeFile)))

//go:embed sqlite/database.go.tmpl
var sqliteDatabaseGoFile []byte

var SqliteDatabaseGoTemplate = template.Must(template.New("sqlite/database.go").Parse(string(sqliteDatabaseGoFile)))

//go:embed sqlite/.env.tmpl
var sqliteEnvFile []byte

var SqliteEnvTemplate = template.Must(template.New("sqlite/.env").Parse(string(sqliteEnvFile)))
