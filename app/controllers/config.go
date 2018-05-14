package controllers

import (
	"net/http"
	"strconv"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/repositories"

	"github.com/labstack/echo"
)

func (co Controllers) GenerateCustomerKey(c echo.Context) error {
	programID := 1

	programCustomerID, err := strconv.Atoi(c.QueryParam("customer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	programCustomerMobile := c.QueryParam("mobile")

	api := new(repositories.API)
	customerKey, err := api.GenerateCustomerKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	c.JSON(http.StatusOK, customerKey)
	return nil
}
