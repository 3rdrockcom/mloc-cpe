package controllers

import (
	"net/http"
	"strconv"

	API "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/api"

	"github.com/labstack/echo"
)

func (co Controllers) GetCustomerKey(c echo.Context) error {
	programID := 1

	programCustomerID, err := strconv.Atoi(c.QueryParam("customer_id"))
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	programCustomerMobile := c.QueryParam("mobile")

	api := API.New()
	customerAccessKey, err := api.GetCustomerAccessKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return SendResponse(c, http.StatusOK, customerAccessKey)

}

func (co Controllers) GenerateCustomerKey(c echo.Context) error {
	programID := 1

	programCustomerID, err := strconv.Atoi(c.QueryParam("customer_id"))
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	programCustomerMobile := c.QueryParam("mobile")

	api := API.New()
	customerAccessKey, err := api.GenerateCustomerAccessKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return SendResponse(c, http.StatusOK, customerAccessKey)
}
