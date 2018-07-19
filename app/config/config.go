package config

import (
	"fmt"
)

// Configuration contains the application configuration
type Configuration struct {
	Application Application
	Path        Path
	Server      Server
	DB          Database
}

// cfg contains the processed configuration values
var cfg Configuration

// New processes the configuration values
func New() (c Configuration, err error) {

	// Application
	c.Application, err = newApplication()
	if err != nil {
		return
	}

	// Path
	c.Path, err = newPath()
	if err != nil {
		return
	}

	// Server
	c.Server, err = newServer()
	if err != nil {
		return
	}

	// Database
	c.DB, err = newDatabase()
	if err != nil {
		return
	}

	// Display application information
	fmt.Print(
		fmt.Sprintf("%-15s %s\n", "Name:", c.Application.Name),
		fmt.Sprintf("%-15s %s\n", "Version:", c.Application.Version),
		fmt.Sprintf("%-15s %s\n", "Build:", c.Application.Build),
		fmt.Sprintf("%-15s %s\n", "Environment:", c.Application.Environment),
	)

	cfg = c
	return
}

// Get gets the processed configuration values
func Get() Configuration {
	return cfg
}
