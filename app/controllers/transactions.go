package controllers

import (
	"strconv"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"

	"github.com/labstack/echo"
)

type payloadTransactions struct {
	Transactions []models.Transaction `json:"transactions" binding:"required"`
}

func (co Controllers) PostAddCustomerTransactions(c echo.Context) error {
	var err error
	db := co.DB.GetInstance()

	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		c.JSON(400, helpers.H{"errors": err.Error()})
		return nil
	}

	payload := payloadTransactions{}
	if err = c.Bind(&payload); err != nil {
		c.JSON(400, helpers.H{"errors": err.Error()})
		return nil
	}

	transactions := payload.Transactions
	for i := range transactions {
		transactions[i].CustomerID = customerID
		transactions[i].DateTime = transactions[i].Date.Time

		if err = transactions[i].Validate(); err != nil {
			c.JSON(400, helpers.H{"errors": err.Error()})
			return nil
		}
	}

	for i := range transactions {
		err = db.Model(&transactions[i]).Insert()
		if err != nil {
			c.JSON(500, helpers.H{"errors": err.Error()})
			return nil
		}
	}

	c.JSON(200, transactions)
	return nil
}
