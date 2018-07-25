package cmd

import (
	"errors"
	"strconv"

	"github.com/epointpayment/mloc-cpe/app/config"
	"github.com/epointpayment/mloc-cpe/app/database/migrations"
	"github.com/epointpayment/mloc-cpe/app/embed"
	"github.com/epointpayment/mloc-cpe/app/log"

	"github.com/spf13/cobra"
)

var migration *migrations.Migration

var migrationDatabase string
var migrationStep int

func init() {
	RootCmd.AddCommand(migrationCommand)
	migrationCommand.PersistentFlags().StringVarP(&migrationDatabase, "database", "d", "default", "The database connection to use.")

	migrationCommand.AddCommand(statusCommand)

	migrationCommand.AddCommand(migrateUpCommand)
	migrateUpCommand.Flags().IntVarP(&migrationStep, "step", "s", 0, "Number of steps to migrate.")

	migrationCommand.AddCommand(migrateDownCommand)
	migrateDownCommand.Flags().IntVarP(&migrationStep, "step", "s", 0, "Number of steps to migrate.")

	migrationCommand.AddCommand(migrateRedoCommand)
	migrationCommand.AddCommand(migrateCreateCommand)
	migrationCommand.AddCommand(migrateForceCommand)

	// Register migrations
	embed.Register("migrations")
}

// migrationCommand manages database migrations
var migrationCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Manage your apps database migrations",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// initiate logging
		log.Start()
		defer log.Stop()

		// load config
		_, err = config.New()
		if err != nil {
			log.Fatal(err)
		}
		log.SetMode(config.Get().Application.Environment)

		// load migrations
		migration, err = migrations.Load(migrationDatabase)
		if err != nil {
			log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		migration.Unload()
	},
}

// statusCommand prints the status of all migrations
var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "Get the current status of database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := migration.Status(); err != nil {
			log.Fatal(err)
		}
	},
}

// migrateUpCommand applies all available migrations
var migrateUpCommand = &cobra.Command{
	Use:   "up",
	Short: "Migrate the database up. By default, it runs all pending migrations.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := migration.Up(migrationStep); err != nil {
			log.Fatal(err)
		}
	},
}

// migrateDownCommand rolls back a single migration from the current version
var migrateDownCommand = &cobra.Command{
	Use:   "down",
	Short: "Migrate the database down. By default, it rolls back the db by 1 migration.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := migration.Down(migrationStep); err != nil {
			log.Fatal(err)
		}
	},
}

// migrateRedoCommand rolls back the most recently applied migration, then runs it again
var migrateRedoCommand = &cobra.Command{
	Use:   "redo",
	Short: "Redo the last migration.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := migration.Redo(); err != nil {
			log.Fatal(err)
		}
	},
}

// migrateCreateCommand creates a new blank migration file
var migrateCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new migration.",
	Run: func(cmd *cobra.Command, args []string) {
		var name string

		if len(args) == 1 {
			name = args[0]
		}

		if err := migration.Create(name); err != nil {
			log.Fatal(err)
		}
	},
}

// migrateForceCommand sets a migration version (ignores dirty state)
var migrateForceCommand = &cobra.Command{
	Use:   "force",
	Short: "Force migration.",
	Run: func(cmd *cobra.Command, args []string) {
		var version int64
		var err error

		if len(args) == 1 {
			version, err = strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				log.Fatal(errors.New("invalid version"))
			}
		} else {
			log.Fatal(errors.New("requires version"))
		}

		if err := migration.Force(int(version)); err != nil {
			log.Fatal(err)
		}
	},
}
