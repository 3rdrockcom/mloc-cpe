package controllers

import (
	"net/http"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/profiler"

	"github.com/jinzhu/now"
	"github.com/labstack/echo"
)

func (co Controllers) GetCustomerProfile(c echo.Context) error {
	var err error

	customerID := c.Get("customerID").(int)

	startDate, err := time.ParseInLocation(
		"20060102",
		c.QueryParam("startDate"), time.UTC)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"error": err.Error()})
		return nil
	}

	endDate, err := time.ParseInLocation(
		"20060102",
		c.QueryParam("endDate"), time.UTC)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"error": err.Error()})
		return nil
	}

	transactions := models.Transactions{}

	sc, err := Customer.New(customerID)
	if err != nil {
		return err
	}

	if transactions, err = sc.Transactions().GetAllByDateRange(startDate, now.New(endDate).EndOfDay()); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.H{"errors": err.Error()})
		return nil
	}

	p := profiler.New(transactions, 2)
	res := p.Run()

	c.JSON(http.StatusOK, res)
	return nil
}
