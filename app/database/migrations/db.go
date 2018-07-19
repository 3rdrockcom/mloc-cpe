package migrations

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Mysql driver for database library
)

// connect opens a connection to a database.
func connect(driver string, dsn string) (*sql.DB, error) {
	// DSN specific changes
	switch driver {
	case "mysql":
		dsn = strings.TrimLeft(dsn, driver+"://")
	}

	// Create database handler
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// disconnect closes a connection to a database.
func disconnect(db *sql.DB) error {
	return db.Close()
}
