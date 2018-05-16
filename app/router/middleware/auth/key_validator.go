package auth

import (
	API "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/api"
	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"

	"github.com/labstack/echo"
)

func DefaultValidator(key string, c echo.Context) (isValid bool, err error) {
	sa := API.New()

	entry, err := sa.GetCustomerKey(key)
	if err != nil {
		return
	}

	if key != entry.Key {
		err = API.ErrInvalidAPIKey
		return
	}

	customerUniqueID := c.QueryParam("cust_unique_id")

	sc, err := Customer.New(entry.CustomerID)
	if err != nil {
		return
	}

	customer, err := sc.Info().Get()
	if err != nil {
		return
	}

	if customer.CustomerUniqueID != customerUniqueID {
		err = Customer.ErrInvalidUniqueCustomerID
		return
	}

	isValid = true
	c.Set("customerID", entry.CustomerID)

	return
}

func RegistrationValidator(key string, c echo.Context) (isValid bool, err error) {
	sa := API.New()

	entry, err := sa.GetRegistrationKey()
	if err != nil {
		return
	}

	if key != entry.Key {
		err = API.ErrInvalidAPIKey
		return
	}

	isValid = true
	return
}

func LoginValidator(key string, c echo.Context) (isValid bool, err error) {
	sa := API.New()

	entry, err := sa.GetLoginKey()
	if err != nil {
		return
	}

	if key != entry.Key {
		err = API.ErrInvalidAPIKey
		return
	}

	isValid = true
	return
}
