package auth

import (
	Customer "github.com/epointpayment/mloc-cpe/app/services/customer"
	"github.com/juju/errors"

	"github.com/labstack/echo"
)

// CUIDValidator is a validator used for customer unique ID auth middleware
func CUIDValidator(customerUniqueID string, c echo.Context) (isValid bool, err error) {
	// Get customer ID from API Key Auth Validator
	customerID := c.Get("_customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Get customer information
	customer, err := sc.Info().Get()
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Check if customer unique ID is a match
	if customer.CustomerUniqueID != customerUniqueID {
		err = errors.Wrap(err, Customer.ErrInvalidUniqueCustomerID)
		return
	}

	// User CUID is valid
	isValid = true

	// Pass user information to context
	c.Set("customerID", customer.ID)
	c.Set("customerUniqueID", customer.CustomerUniqueID)
	return
}
