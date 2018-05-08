package controllers

import (
	"net/http"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/repositories"

	"github.com/labstack/echo"
)

func (co Controllers) PostAddCustomer(c echo.Context) error {
	var err error
	customer := new(models.Customer)

	if err = c.Bind(customer); err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	if err = customer.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	rc := new(repositories.Customers)
	if err = rc.Create(customer); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.H{"errors": err.Error()})
		return nil
	}

	c.JSON(http.StatusOK, customer)
	return nil
}
