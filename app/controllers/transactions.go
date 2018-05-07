package controllers

import (
	"strconv"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	"github.com/gin-gonic/gin"
)

type payloadTransactions struct {
	Transactions []models.Transaction `json:"transactions" binding:"required"`
}

func (co Controllers) PostAddCustomerTransactions(c *gin.Context) {
	var err error
	db := co.DB.GetInstance()

	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		c.JSON(400, gin.H{"errors": err.Error()})
		return
	}

	payload := payloadTransactions{}
	if err = c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"errors": err.Error()})
		return
	}

	transactions := payload.Transactions
	for i := range transactions {
		transactions[i].CustomerID = customerID
		transactions[i].DateTime = transactions[i].Date.Time

		if err = transactions[i].Validate(); err != nil {
			c.JSON(400, gin.H{"errors": err.Error()})
			return
		}
	}

	for i := range transactions {
		err = db.Model(&transactions[i]).Insert()
		if err != nil {
			c.JSON(500, gin.H{"errors": err.Error()})
			return
		}
	}

	c.JSON(200, transactions)
}
