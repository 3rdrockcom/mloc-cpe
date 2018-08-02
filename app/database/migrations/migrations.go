package migrations

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/epointpayment/mloc-cpe/app/config"
	"github.com/epointpayment/mloc-cpe/app/embed"
	"github.com/epointpayment/mloc-cpe/app/log"

	packrdriver "github.com/fiskeben/packr-source-driver/driver"
	"github.com/gobuffalo/packr"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

func init() {
	// Set migrations table
	mysql.DefaultMigrationsTable = "migrations"
}

// numLeadingZeros is the max number of filler digits for migrations (000001, 000002, 000003, ...)
const numLeadingZeros = 6

// Migration contains required information for performing migration operations
type Migration struct {
	database string
	driver   string
	dsn      string
	path     string
	ext      string
	box      packr.Box
	migrate  *migrate.Migrate
}

// Load will load migration and database information.
func Load(database string) (m *Migration, err error) {
	migrationsDir := config.Get().Path.Migrations
	m = &Migration{
		database: database,
		path:     migrationsDir + string(os.PathSeparator) + database,
	}

	// Check if database name is not empty
	if m.database == "" {
		err = errors.New("No database selected")
		return
	}

	// Database found
	if m.database == "default" {
		entry := config.Get().DB
		m.driver = entry.Driver
		m.dsn = entry.DSN
	}

	if m.driver == "" || m.dsn == "" {
		err = errors.New("database not found")
		return
	}

	// Check available drivers
	switch m.driver {
	case "mysql":
		m.ext = "sql"
	default:
		err = fmt.Errorf("%q database driver not supported", m.driver)
		return
	}

	// Get box
	m.box, err = embed.Get(m.path)
	if err != nil {
		err = fmt.Errorf("box has not been configured for %s", m.database)
		return
	}

	// Setup migration driver
	migrateDriver, err := packrdriver.WithInstance(m.box)
	if err != nil {
		err = fmt.Errorf("failed to create migration data driver: %v", err)
		return
	}

	// Setup migration source
	m.migrate, err = migrate.NewWithSourceInstance("packr", migrateDriver, m.dsn)
	if err != nil {
		return
	}

	// Setup logger
	m.migrate.Log = new(Logger)

	return
}

// Unload will unload migration and database information.
func (m *Migration) Unload() (err error) {
	if errSource, errDatabase := m.migrate.Close(); errSource != nil || errDatabase != nil {
		return err
	}
	return
}

// CreateFolder will set the proper migrations folder.
func (m *Migration) CreateFolder() (err error) {
	log.Printf("migrations: creating migrations folder for currently selected")
	err = os.MkdirAll(m.path, os.ModePerm)
	return
}

//Status prints the status of all migrations.
func (m *Migration) Status() (err error) {
	version, isDirty, err := m.migrate.Version()
	if err != nil {
		return
	}

	// Display migration information
	fmt.Print(
		fmt.Sprintf("\nMIGRATION STATUS\n"),
		fmt.Sprintf("---\n"),
		fmt.Sprintf("%-15s %v\n", "Version:", version),
		fmt.Sprintf("%-15s %v\n", "Dirty:", isDirty),
	)

	return
}

// Up applies all available migrations.
func (m *Migration) Up(step int) (err error) {
	// Upgrade to specific migration
	if step > 0 {
		log.Printf("migrations: upgrading %d migration(s)", step)
		if err := m.migrate.Steps(step); err != nil && err != migrate.ErrNoChange {
			return err
		}

		return
	}

	// Upgrade all migrations
	log.Printf("migrations: upgrading to latest available migration")
	if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return
}

// Down rolls back a single migration from the current version.
func (m *Migration) Down(step int) (err error) {
	// Downgrade to specific migration
	if step > 0 {
		log.Printf("migrations: downgrading %d migration(s)", step)
		if err := m.migrate.Steps(-1 * step); err != nil && err != migrate.ErrNoChange {
			return err
		}

		return
	}

	// Downgrade all migrations
	log.Printf("migrations: downgrading to oldest available migration")
	if err := m.migrate.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return
}

// Redo rolls back the most recently applied migration, then runs it again.
func (m *Migration) Redo() error {
	// Downgrade migration
	log.Println("migrations: redo migration - step down [1/2]")
	if err := m.Down(1); err != nil {
		return err
	}

	// Upgrade migration
	log.Println("migrations: redo migration - step up [2/2]")
	if err := m.Up(1); err != nil {
		return err
	}

	return nil
}

// Create writes a new blank migration file.
func (m *Migration) Create(name string) (err error) {
	// Check if name is valid
	if name == "" {
		err = errors.New("Please specify a name for the migration")
		return
	}

	// Create migrations folder for database
	err = m.CreateFolder()
	if err != nil {
		return
	}

	// Determine the next version number
	version, err := m.getNextVersion(m.box.List())
	if err != nil {
		return
	}

	// Create migration files
	for _, direction := range []string{"up", "down"} {
		// Determine filename
		filename := fmt.Sprintf("%v%v_%v.%v.%v", m.path+string(os.PathSeparator), version, name, direction, m.ext)

		// Create file
		_, err := os.Create(filename)
		if err != nil {
			return err
		}

		log.Println("migrations: created migration file - " + filename)
	}

	return
}

// Force sets a migration version.
func (m *Migration) Force(version int) (err error) {
	err = m.migrate.Force(version)
	if err != nil {
		return
	}

	log.Println("migrations: set migration version to " + strconv.FormatInt(int64(version), 10))
	return
}

func (m *Migration) getNextVersion(filenames []string) (nextVersion string, err error) {

	// Check if digit padding is a positive number
	if numLeadingZeros <= 0 {
		err = errors.New("Number of leading zeros must be a positive number")
		return
	}

	// Determine next version increment
	version := 1
	if len(filenames) > 0 {

		// Setup migration format pattern
		r, err := regexp.Compile(`^(\d{` + strconv.FormatInt(numLeadingZeros, 10) + `})\_(\S+)\.(\S+)\.` + m.ext + `$`)
		if err != nil {
			return "", err
		}

		// Check if filenames are valid and add to list
		migrations := []string{}
		for i, filename := range filenames {
			// Ignore non-migration files
			if !strings.HasSuffix(filename, m.ext) {
				continue
			}

			// Check if filename matches migration format pattern
			filename = strings.TrimPrefix(filename, m.path)
			match := r.MatchString(filename)
			if !match {
				err = errors.New("Malformed migration filename: " + filename)
				return "", err
			}

			// Add to migrations file list
			migrations = append(migrations, filenames[i])
		}

		// Sort list of migrations
		sort.Strings(migrations)

		// Get last file name in list (most recent migration)
		migration := migrations[len(migrations)-1]

		// Check if migration filename is valid
		versionStr := r.FindStringSubmatch(migration)[1]

		// Convert the version string to an integer
		version, err = strconv.Atoi(versionStr)
		if err != nil {
			return "", err
		}

		// Increment the current version to the next version
		version++
	}
	if version <= 0 {
		err = errors.New("Next sequence number must be positive")
		return
	}

	// Check if the next version will hit sequence limit (leading zero)
	nextVersion = strconv.Itoa(version)
	if len(nextVersion) > numLeadingZeros {
		return "", fmt.Errorf("Next sequence number %s too large. At most %d digits are allowed", nextVersion, numLeadingZeros)
	}

	// Prepend zeros
	padding := numLeadingZeros - len(nextVersion)
	if padding > 0 {
		nextVersion = strings.Repeat("0", padding) + nextVersion
	}

	return
}
