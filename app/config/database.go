package config

import (
	"strconv"

	"github.com/go-ozzo/ozzo-validation/is"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/spf13/viper"
)

// Database contains database information
type Database struct {
	Driver   string
	Host     string
	Port     int64
	Database string
	Username string
	Password string
	Flags    string
	DSN      string
}

// Validate checks the configuration for invalid values
func (c Database) Validate() (err error) {
	err = validation.Errors{
		"DB_CONNECTION": validation.Validate(c.Driver, validation.Required, validation.In("mysql")),
		"DB_HOST":       validation.Validate(c.Host, validation.Required, is.Host),
		"DB_PORT":       validation.Validate(strconv.FormatInt(c.Port, 10), validation.Required, is.Port),
		"DB_DATABASE":   validation.Validate(c.Database, validation.Required),
		"DB_USERNAME":   validation.Validate(c.Username, validation.Required),
		"DB_PASSWORD":   validation.Validate(c.Password, validation.Required),
		"DB_DSN":        validation.Validate(c.DSN, validation.Required, is.RequestURL),
	}.Filter()

	return
}

// newDatabase configures database settings
func newDatabase() (c Database, err error) {
	// Bind
	viper.BindEnv("DB_CONNECTION")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_DATABASE")
	viper.BindEnv("DB_USERNAME")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_FLAGS")

	// Set defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 3306)

	// Get values and assign
	c = Database{
		Driver:   viper.GetString("DB_CONNECTION"),
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetInt64("DB_PORT"),
		Database: viper.GetString("DB_DATABASE"),
		Username: viper.GetString("DB_USERNAME"),
		Password: viper.GetString("DB_PASSWORD"),
		Flags:    viper.GetString("DB_FLAGS"),
	}
	c.DSN = generateDSN(c)

	// Validate
	err = c.Validate()
	if err != nil {
		return
	}

	return
}

// generateDSN creates a DSN from database config
func generateDSN(d Database) string {
	return d.Driver + "://" + d.Username + ":" + d.Password + "@tcp(" + d.Host + ":" + strconv.FormatInt(d.Port, 10) + ")/" + d.Database + "?" + d.Flags
}
