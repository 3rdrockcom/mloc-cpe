package controllers

import (
	"net/http"

	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"

	"github.com/labstack/echo"
)

func (co Controllers) PostAddCustomer(c echo.Context) error {
	customerID := c.Get("customerID").(int)

	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	customer, err := sc.Info().Get()
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err = c.Bind(customer); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err = customer.Validate(); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err = sc.Info().Update(customer); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return SendResponse(c, http.StatusOK, customer)
}
