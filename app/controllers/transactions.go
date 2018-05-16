package controllers

import (
	"net/http"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
)

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
