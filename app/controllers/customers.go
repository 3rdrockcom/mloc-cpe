package controllers

import (
	"net/http"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"
	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"

	"github.com/labstack/echo"
)

func (co Controllers) PostAddCustomer(c echo.Context) error {
	var err error

	customerID := c.Get("customerID").(int)

	sc, err := Customer.New(customerID)
	if err != nil {
		return err
	}
	customer, err := sc.Info().Get()
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	if err = c.Bind(customer); err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	if err = customer.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	if err = sc.Info().Update(customer); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.H{"errors": err.Error()})
		return nil
	}

	c.JSON(http.StatusOK, customer)
	return nil
}
