package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

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
		Migrations: "app/migrations",
	}

	// Validate
	err = c.Validate()
	if err != nil {
		return
	}

	return
}
