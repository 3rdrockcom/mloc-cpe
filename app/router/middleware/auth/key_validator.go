package auth

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/repositories"

	"github.com/labstack/echo"
)

func DefaultValidator(key string, c echo.Context) (bool, error) {
	var err error
	var customerUniqueID string

	api := new(repositories.API)
	customerKey, err := api.GetAPICustomerKey(key)
	if err != nil {
		return false, err
	}

	customerUniqueID = c.QueryParam("cust_unique_id")

	customers := new(repositories.Customers)
	_, err = customers.GetByCustomerUniqueID(customerUniqueID)
	if err != nil {
		return false, err
	}
	//verify

	c.Set("customerID", customerKey.CustomerID)

	return true, nil
}

func RegistrationValidator(key string, c echo.Context) (bool, error) {
	var err error

	api := new(repositories.API)
	registrationKey, err := api.GetAPIRegistrationKey()
	if err != nil {
		return false, err
	}

	if key != registrationKey.Key {
		return false, nil
	}

	return true, nil
}

func LoginValidator(key string, c echo.Context) (bool, error) {
	var err error

	api := new(repositories.API)
	loginKey, err := api.GetAPILoginKey()
	if err != nil {
		return false, err
	}

	if key != loginKey.Key {
		return false, nil
	}

	return true, nil
}
