package controllers

import (
	"net/http"
	"strconv"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"
	API "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/api"

	"github.com/labstack/echo"
)

func (co Controllers) GetCustomerKey(c echo.Context) error {
	programID := 1

	programCustomerID, err := strconv.Atoi(c.QueryParam("customer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	programCustomerMobile := c.QueryParam("mobile")

	api := API.New()
	customerAccessKey, err := api.GetCustomerAccessKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	c.JSON(http.StatusOK, customerAccessKey)
	return nil

}
