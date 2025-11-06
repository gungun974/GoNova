package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5"     // Import postgres driver
	_ "github.com/mattn/go-sqlite3" // Import sqlite3 driver

	"github.com/gungun974/gonova/internal/logger"
	"github.com/gungun974/gonova/internal/utils"
)

func getMigrationSourcePath() string {
	projectPath := "."

	_, err := utils.GetGoModName(projectPath)
	if err != nil {
		logger.MainLogger.Fatalf("Can't parse go mod : %v", err)
	}

	return filepath.Join(projectPath, "/internal/database/migrations")
}

func getMigrateInstance() *migrate.Migrate {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		logger.DatabaseLogger.Fatalln("DATABASE_URL is not set")
	}

	m, err := migrate.New(fmt.Sprintf("file://%s", getMigrationSourcePath()), databaseURL)
	if err != nil {
		logger.DatabaseLogger.Fatalf("Unable to get migrate instance: %v", err)
	}

	return m
}

func MigrateCurrent() {
	m := getMigrateInstance()

	version, dirty, err := m.Version()

	if errors.Is(err, migrate.ErrNilVersion) {
		logger.DatabaseLogger.Info("No current migration are applied")
		return
	}

	if err != nil {
		logger.DatabaseLogger.Fatalf("Migration version failed: %v", err)
	}

	if dirty {
		logger.DatabaseLogger.Warnf("Current migration is %v and dirty", version)
		return
	}
	logger.DatabaseLogger.Infof("Current migration is %v", version)
}

func MigrateUp() {
	m := getMigrateInstance()

	logger.DatabaseLogger.Info("Start migrations up")

	start := time.Now()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.DatabaseLogger.Fatalf("Migration up failed: %v", err)
	}

	duration := time.Since(start)

	logger.DatabaseLogger.Infof("Finish making migrations in %s", duration)

	MigrateCurrent()
}

func MigrateDown() {
	m := getMigrateInstance()

	logger.DatabaseLogger.Info("Start migrations down")

	start := time.Now()

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.DatabaseLogger.Fatalf("Migration down failed: %v", err)
	}

	duration := time.Since(start)

	logger.DatabaseLogger.Infof("Finish making migrations in %s", duration)

	MigrateCurrent()
}

func MigrateVersion(version int) {
	m := getMigrateInstance()

	logger.DatabaseLogger.Info("Start migration force")

	start := time.Now()

	if err := m.Force(version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.DatabaseLogger.Fatalf("Migration force failed: %v", err)
	}

	duration := time.Since(start)

	logger.DatabaseLogger.Infof("Finish forcing migration in %s", duration)

	MigrateCurrent()
}

func MigrateDrop() {
	m := getMigrateInstance()

	logger.DatabaseLogger.Info("Start migration drop")

	start := time.Now()

	if err := m.Drop(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.DatabaseLogger.Fatalf("Migration drop failed: %v", err)
	}

	duration := time.Since(start)

	logger.DatabaseLogger.Infof("Finish droping migration in %s", duration)

	MigrateCurrent()
}

type Migration struct {
	Version int
	Name    string
}

func ListMigrations() []Migration {
	dir := filepath.Clean(getMigrationSourcePath())

	matches, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		log.Fatal(err)
	}

	var migrations []Migration
	re := regexp.MustCompile(`^(\d+)_.*\.sql$`)

	for _, path := range matches {
		filename := filepath.Base(path)

		if re.MatchString(filename) {
			match := re.FindStringSubmatch(filename)
			versionStr := match[1]
			version, err := strconv.Atoi(versionStr)
			if err != nil {
				continue
			}

			migrations = append(migrations, Migration{
				Version: version,
				Name:    strings.TrimSuffix(filename, ".up.sql"),
			})
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations
}

func MigrateCreate(name string) {
	dir := filepath.Clean(getMigrationSourcePath())

	version := time.Now().UTC().Format("20060102150405")

	versionGlob := filepath.Join(dir, version+"_*.sql")
	matches, err := filepath.Glob(versionGlob)
	if err != nil {
		logger.DatabaseLogger.Fatal(err)
	}

	if len(matches) > 0 {
		logger.DatabaseLogger.Fatalf("Duplicate migration version: %s", version)
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		logger.DatabaseLogger.Fatal(err)
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s.sql", version, name, direction)
		filename := filepath.Join(dir, basename)

		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			logger.DatabaseLogger.Fatal(err)
		}

		_ = f.Close()

		absPath, _ := filepath.Abs(filename)
		logger.DatabaseLogger.Infof("Create file: %s", absPath)
	}

	emptyMigrationFilePath := filepath.Join(dir, "EMPTY_MIGRATION.sql")

	if _, err := os.Stat(emptyMigrationFilePath); err == nil {
		_ = os.Remove(emptyMigrationFilePath)
	}
}
