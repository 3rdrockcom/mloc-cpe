package config

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/viper"
)

// Server contains server information
type Server struct {
	Host string
	Port int64
}

// Validate checks the configuration for invalid values
func (c Server) Validate() (err error) {
	err = validation.Errors{
		"HOST": validation.Validate(c.Host, validation.Required, is.Host),
		"PORT": validation.Validate(strconv.FormatInt(c.Port, 10), validation.Required, is.Port),
	}.Filter()

	return
}

// newServer configures server settings
func newServer() (c Server, err error) {
	// Bind
	viper.BindEnv("HOST")
	viper.BindEnv("PORT")

	// Set defaults
	viper.SetDefault("HOST", "localhost")
	viper.SetDefault("PORT", 3000)

	// Get values and assign
	c = Server{
		Host: viper.GetString("HOST"),
		Port: viper.GetInt64("PORT"),
	}

	// Validate
	err = c.Validate()
	if err != nil {
		return
	}

	return
}
