package controllers

import (
	"net/http"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/profiler"

	"github.com/jinzhu/now"
	"github.com/labstack/echo"
)

func (co Controllers) GetCustomerProfile(c echo.Context) error {
	customerID := c.Get("customerID").(int)

	startDate, err := time.ParseInLocation(
		"20060102",
		c.QueryParam("startDate"), time.UTC)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	endDate, err := time.ParseInLocation(
		"20060102",
		c.QueryParam("endDate"), time.UTC)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	sc, err := Customer.New(customerID)
	if err != nil {
		return err
	}

	transactions := models.Transactions{}
	if transactions, err = sc.Transactions().GetAllByDateRange(startDate, now.New(endDate).EndOfDay()); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	p := profiler.New(transactions, 2)
	res := p.Run()

	return SendResponse(c, http.StatusOK, res)
}
