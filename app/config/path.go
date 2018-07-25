package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// PathMigrations is the location where database migrations are located
const PathMigrations = "app/migrations"

// Path contains path information
type Path struct {
	Migrations string
}

// Validate checks the configuration for invalid values
func (c Path) Validate() (err error) {
	err = validation.ValidateStruct(&c,
		validation.Field(&c.Migrations, validation.Required),
	)

	return
}

// newServer configures server settings
func newPath() (c Path, err error) {
	// Get values and assign
	c = Path{
		Migrations: PathMigrations,
	}

	// Validate
	err = c.Validate()
	if err != nil {
		return
	}

	return
}
