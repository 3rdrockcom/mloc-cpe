package controllers

import (
	"net/http"
	"strconv"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/repositories"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
)

type payloadTransactions struct {
	Transactions []models.Transaction `json:"transactions" binding:"required"`
}

func (co Controllers) PostAddCustomerTransactions(c echo.Context) error {
	var err error

	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	payload := payloadTransactions{}
	if err = c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	transactions := payload.Transactions
	for i := range transactions {
		transactions[i].CustomerID = customerID
		transactions[i].DateTime = transactions[i].Date.Time
	}

	if err = validation.Validate(transactions); err != nil {
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

	rt := new(repositories.Transactions)
	if err = rt.Create(customerID, transactions); err != nil {
		c.JSON(http.StatusInternalServerError, helpers.H{"errors": err.Error()})
		return nil
	}

	c.JSON(http.StatusOK, transactions)
	return nil
}
