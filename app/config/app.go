package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/spf13/viper"
)

const (
	// EnvDevelopment is the setting for a development environment
	EnvDevelopment = "development"

	// EnvProduction is the setting for a production environment
	EnvProduction = "production"
)

var (
	// Version is the application semantic version
	Version string

	// Build is the application build version
	Build string
)

func init() {
	if Version == "" {
		Version = "unknown"
	}
	if Build == "" {
		Build = "unknown"
	}
}

// Application contains application information
type Application struct {
	Name        string
	Build       string
	Version     string
	Environment string
}

// Validate checks the configuration for invalid values
func (c Application) Validate() (err error) {
	err = validation.Errors{
		"NAME":        validation.Validate(c.Name, validation.Required),
		"ENVIRONMENT": validation.Validate(c.Environment, validation.Required, validation.In(EnvDevelopment, EnvProduction)),
	}.Filter()

	return
}

// newApplication configures application settings
func newApplication() (c Application, err error) {
	// Bind
	viper.BindEnv("NAME")
	viper.BindEnv("ENVIRONMENT")

	// Set defaults
	viper.SetDefault("NAME", "app")
	viper.SetDefault("ENVIRONMENT", EnvProduction)

	// Get values and assign
	c = Application{
		Name:        viper.GetString("NAME"),
		Build:       Build,
		Version:     Version,
		Environment: viper.GetString("ENVIRONMENT"),
	}

	// Validate
	err = c.Validate()
	if err != nil {
		return
	}

	return
}

// IsDev determines if the application environment is in development mode
func IsDev() bool {
	return cfg.Application.Environment == EnvDevelopment
}

// IsProd determines if the application environment is in production mode
func IsProd() bool {
	return cfg.Application.Environment == EnvProduction
}
