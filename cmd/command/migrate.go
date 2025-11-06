package command

import (
	"strconv"

	"github.com/gungun974/gonova/internal/database"
	"github.com/gungun974/gonova/internal/form"
	"github.com/gungun974/gonova/internal/helpers"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func init() {
	databaseCmd := &cobra.Command{Use: "migrate"}

	databaseCmd.AddCommand(migrateCreateCmd)

	databaseCmd.AddCommand(migrateVersionCmd)

	databaseCmd.AddCommand(migrateUpCmd)

	databaseCmd.AddCommand(migrateDownCmd)

	migrateDownCmd.Flags().BoolP("force", "f", false, "Force down migration")

	databaseCmd.AddCommand(migrateForceCmd)

	migrateForceCmd.Flags().BoolP("force", "f", false, "Force version migration")

	databaseCmd.AddCommand(migrateDropCmd)

	migrateDropCmd.Flags().BoolP("force", "f", false, "Force drop database")

	rootCmd.AddCommand(databaseCmd)
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create new timestamped UTC up and down migrations",
	Args:  cobra.MinimumNArgs(0),
	Run:   MigrateCreate,
}

func MigrateCreate(cmd *cobra.Command, args []string) {
	var migrationName string

	if len(args) == 0 {
		migrationName = form.AskInputWithPlaceholder(
			"Migration name :",
			"Create posts table",
		)
	} else {
		migrationName = args[0]
	}

	if err := godotenv.Load(); err == nil {
		logger.MainLogger.Info("Loading .env file")
	}

	database.MigrateCreate(helpers.ToSnakeCase(migrationName))
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get database migration version",
	Run:   MigrateVersion,
}

func MigrateVersion(_ *cobra.Command, _ []string) {
	err := godotenv.Load()

	if err == nil {
		logger.MainLogger.Info("Loading .env file")
	}

	database.MigrateCurrent()
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Execute every database migrations up",
	Run:   MigrateUp,
}

func MigrateUp(_ *cobra.Command, _ []string) {
	err := godotenv.Load()

	if err == nil {
		logger.MainLogger.Info("Loading .env file")
	}

	database.MigrateUp()
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Execute database migrations down (1 step)",
	Run:   MigrateDown,
}

func MigrateDown(cmd *cobra.Command, _ []string) {
	force, _ := cmd.Flags().GetBool("force")

	if !force {
		logger.MainLogger.Fatal("Migration down can destroy data, use \"--force\" if you are sure")
	}

	err := godotenv.Load()

	if err == nil {
		logger.MainLogger.Info("Loading .env file")
	}

	database.MigrateDown()
}

var migrateForceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Change database migration version by force",
	Args:  cobra.MinimumNArgs(0),
	Run:   MigrateForce,
}

func MigrateForce(cmd *cobra.Command, args []string) {
	force, _ := cmd.Flags().GetBool("force")

	if !force {
		logger.MainLogger.Fatal(
			"Migration force can mess with migration history, use \"--force\" if you are sure",
		)
	}

	var version int

	if len(args) == 0 {
		migrations := database.ListMigrations()

		choices := []form.Choice[int]{}

		for _, migration := range migrations {
			choices = append(choices, form.Choice[int]{
				Name:  migration.Name,
				Value: migration.Version,
			})
		}

		version = form.AskChoiceSearch(
			"Migration :",
			choices,
		)
	} else {
		var err error
		version, err = strconv.Atoi(args[0])
		if err != nil {
			logger.MainLogger.Fatal("Version is not a valid int")
		}
	}

	if err := godotenv.Load(); err == nil {
		logger.MainLogger.Info("Loading .env file")
	}

	database.MigrateVersion(version)
}

var migrateDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop everything inside database",
	Run:   MigrateDrop,
}

func MigrateDrop(cmd *cobra.Command, args []string) {
	force, _ := cmd.Flags().GetBool("force")

	if !force {
		logger.MainLogger.Fatal(
			"Migration drop will destroy everything inside the database, use \"--force\" if you are sure",
		)
	}

	err := godotenv.Load()

	if err == nil {
		logger.MainLogger.Info("Loading .env file")
	}

	database.MigrateDrop()
}
