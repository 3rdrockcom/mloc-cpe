package auth

import (
	API "github.com/epointpayment/mloc-cpe/app/services/api"
	"github.com/juju/errors"

	"github.com/labstack/echo"
)

// CustomerValidator is a validator used for customer key auth middleware
func CustomerValidator(key string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Get customer key
	entry, err := sa.GetCustomerKey(key)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Check if key is a match
	if key != entry.Key {
		err = errors.Wrap(err, API.ErrInvalidAPIKey)
		return
	}

	// User key is valid
	isValid = true

	// Pass user information to context
	c.Set("_customerID", entry.CustomerID)
	return
}

// RegistrationValidator is a validator used for registration key auth middleware
func RegistrationValidator(key string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Get API key for registration
	entry, err := sa.GetRegistrationKey()
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Check if key is a match
	if key != entry.Key {
		err = errors.Wrap(err, API.ErrInvalidAPIKey)
		return
	}

	// User is authorized
	isValid = true
	return
}

// LoginValidator is a validator used for login key auth middleware
func LoginValidator(key string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Get API key for login
	entry, err := sa.GetLoginKey()
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Check if key is a match
	if key != entry.Key {
		err = errors.Wrap(err, API.ErrInvalidAPIKey)
		return
	}

	// User is authorized
	isValid = true
	return
}
