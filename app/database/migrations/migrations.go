package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/epointpayment/mloc-cpe/app/config"
	"github.com/epointpayment/mloc-cpe/app/log"

	"github.com/pressly/goose"
)

// Migration struct.
type Migration struct {
	db   *sql.DB
	path string
}

// Load will load migration and database information.
func Load(database string) (*Migration, error) {
	var migration = &Migration{
		path: config.Get().Path.Migrations + string(os.PathSeparator) + database,
	}
	var err error

	// check if database name is not empty
	if database == "" {
		return nil, errors.New("No database selected")
	}

	// get database configuration
	var driver, dsn string

	// database found
	if database == "default" {
		entry := config.Get().DB
		driver = entry.Driver
		dsn = entry.DSN
	}

	if driver == "" || dsn == "" {
		return nil, errors.New("database not found")
	}

	if err := initGoose(driver); err != nil {
		log.Fatalln(err)
	}

	// connect to database
	migration.db, err = connect(driver, dsn)
	if err != nil {
		return nil, err
	}

	// create migrations folder
	err = migration.CreateFolder()
	if err != nil {
		return nil, errors.New("unable to validate migrations folder: " + migration.path)
	}

	return migration, nil
}

// initGoose prepares Goose migrations tool for use
func initGoose(driver string) (err error) {
	// select database dialect
	switch driver {
	case "postgres", "mysql", "sqlite3", "redshift":
		if err := goose.SetDialect(driver); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%q driver not supported", driver)
	}

	// Setup logger
	goose.SetLogger(log.DefaultLogger)

	return
}

// Unload will unload migration and database information.
func (m *Migration) Unload() error {
	err := disconnect(m.db)
	return err
}

// CreateFolder will set the proper migrations folder.
func (m *Migration) CreateFolder() error {
	err := os.MkdirAll(m.path, 0766)
	return err
}

//Status prints the status of all migrations.
func (m *Migration) Status() error {
	if err := goose.Run("status", m.db, m.path); err != nil {
		return err
	}
	return nil
}

// Up applies all available migrations.
func (m *Migration) Up(step int) error {

	if step > 0 {
		migrations, err := m.getMigrations()
		if err != nil {
			return err
		}

		for i := 0; i < step; i++ {
			currentVersion, err := m.getCurrentVersion()
			if err != nil {
				return err
			}

			next, err := migrations.Next(currentVersion)
			if err != nil {
				if err == goose.ErrNoNextVersion {
					return fmt.Errorf("no migration %v\n", currentVersion)
				}
				return err
			}

			if err = next.Up(m.db); err != nil {
				return err
			}
		}
	} else {
		if err := goose.Run("up", m.db, m.path); err != nil {
			return err
		}
	}

	return nil
}

// Down rolls back a single migration from the current version.
func (m *Migration) Down(step int) error {
	if step > 0 {
		migrations, err := m.getMigrations()
		if err != nil {
			return err
		}

		for i := 0; i < step; i++ {
			currentVersion, err := m.getCurrentVersion()
			if err != nil {
				return err
			}

			current, err := migrations.Current(currentVersion)
			if err != nil {
				return fmt.Errorf("no migration %v\n", currentVersion)
			}

			if err = current.Down(m.db); err != nil {
				return err
			}
		}
	} else {
		if err := goose.Run("down", m.db, m.path); err != nil {
			return err
		}
	}

	return nil
}

// Redo rolls back the most recently applied migration, then runs it again.
func (m *Migration) Redo() error {
	if err := goose.Run("redo", m.db, m.path); err != nil {
		return err
	}

	return nil
}

// Create writes a new blank migration file.
func (m *Migration) Create(name string) error {
	if name == "" {
		return errors.New("Please specify a name for the migration")
	}

	if err := goose.Create(m.db, m.path, name, "go"); err != nil {
		return err
	}

	return nil
}

// getCurrentVersion retrieves the current version for this database
func (m *Migration) getCurrentVersion() (int64, error) {
	currentVersion, err := goose.GetDBVersion(m.db)
	if err != nil {
		return 0, err
	}

	return currentVersion, nil
}

// getMigrations returns all the valid migration scripts in the migrations folder
func (m *Migration) getMigrations() (goose.Migrations, error) {
	migrations, err := goose.CollectMigrations(m.path, int64(0), int64((1<<63)-1))
	if err != nil {
		return nil, err
	}

	return migrations, nil
}
