package controllers

import (
	"net/http"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/profiler"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jinzhu/now"
	"github.com/labstack/echo"
)

func (co Controllers) GetCustomer(c echo.Context) error {
	customerID := c.Get("customerID").(int)

	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	customer, err := sc.Info().GetDetails()
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return SendResponse(c, http.StatusOK, customer)
}

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

type payloadTransactions struct {
	Transactions []models.Transaction `json:"transactions" binding:"required"`
}

func (co Controllers) PostAddCustomerTransactions(c echo.Context) error {
	customerID := c.Get("customerID").(int)

	payload := payloadTransactions{}
	if err := c.Bind(&payload); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	transactions := payload.Transactions
	for i := range transactions {
		transactions[i].CustomerID = customerID
		transactions[i].DateTime = transactions[i].Date.Time
	}

	if err := validation.Validate(transactions); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	sc, err := Customer.New(customerID)
	if err != nil {
		return err
	}

	if err = sc.Transactions().Create(transactions); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return SendResponse(c, http.StatusOK, transactions)
}

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
