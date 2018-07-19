package controllers

import (
	"net/http"
	"strconv"

	API "github.com/epointpayment/mloc-cpe/app/services/api"

	"github.com/juju/errors"
	"github.com/labstack/echo"
)

// GetCustomerKey retrieves an API key from the database
func (co Controllers) GetCustomerKey(c echo.Context) (err error) {
	// Get program ID
	programID, err := strconv.Atoi(c.QueryParam("program_id"))
	if err != nil {
		err = errors.Wrap(err, API.ErrInvalidProgramID)
		return
	}

	// Get program customer ID
	programCustomerID, err := strconv.Atoi(c.QueryParam("program_customer_id"))
	if err != nil {
		err = errors.Wrap(err, API.ErrInvalidProgramCustomerID)
		return
	}

	// Get program customer mobile
	programCustomerMobile := c.QueryParam("program_customer_mobile")

	// Get API key
	api := API.New()
	customerAccessKey, err := api.GetCustomerAccessKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Send results
	return SendResponse(c, http.StatusOK, customerAccessKey)

}

// GenerateCustomerKey creates a new customer and API key and stores it in the database
func (co Controllers) GenerateCustomerKey(c echo.Context) (err error) {
	// Get program ID
	programID, err := strconv.Atoi(c.QueryParam("program_id"))
	if err != nil {
		err = errors.Wrap(err, API.ErrInvalidProgramID)
		return
	}

	// Get program customer ID
	programCustomerID, err := strconv.Atoi(c.QueryParam("program_customer_id"))
	if err != nil {
		err = errors.Wrap(err, API.ErrInvalidProgramCustomerID)
		return
	}

	// Get program customer mobile
	programCustomerMobile := c.QueryParam("program_customer_mobile")

	// Get API key
	api := API.New()
	customerAccessKey, err := api.GenerateCustomerAccessKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Send results
	return SendResponse(c, http.StatusOK, customerAccessKey)
}
