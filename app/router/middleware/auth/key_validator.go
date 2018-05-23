package auth

import (
	API "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/api"
	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"

	"github.com/labstack/echo"
)

func DefaultValidator(key string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Get customer key
	entry, err := sa.GetCustomerKey(key)
	if err != nil {
		return
	}

	// Check if key is a match
	if key != entry.Key {
		err = API.ErrInvalidAPIKey
		return
	}

	// Get customer unique ID
	customerUniqueID := c.QueryParam("cust_unique_id")

	// Initialize customer service
	sc, err := Customer.New(entry.CustomerID)
	if err != nil {
		return
	}

	// Get customer information
	customer, err := sc.Info().Get()
	if err != nil {
		return
	}

	// Check if customer unique ID is a match
	if customer.CustomerUniqueID != customerUniqueID {
		err = Customer.ErrInvalidUniqueCustomerID
		return
	}

	// User is authorized
	isValid = true

	// Pass user information to context
	c.Set("customerID", entry.CustomerID)
	return
}

func RegistrationValidator(key string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Get API key for registration
	entry, err := sa.GetRegistrationKey()
	if err != nil {
		return
	}

	// Check if key is a match
	if key != entry.Key {
		err = API.ErrInvalidAPIKey
		return
	}

	// User is authorized
	isValid = true
	return
}

// LoginValidator is a validator used for key auth middleware
func LoginValidator(key string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Get API key for login
	entry, err := sa.GetLoginKey()
	if err != nil {
		return
	}

	// Check if key is a match
	if key != entry.Key {
		err = API.ErrInvalidAPIKey
		return
	}

	// User is authorized
	isValid = true
	return
}
