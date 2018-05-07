package controllers

import (
	"strconv"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/profiler"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

func (co Controllers) GetCustomerProfile(c *gin.Context) {
	var err error
	db := co.DB.GetInstance()

	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		c.JSON(400, gin.H{"errors": err.Error()})
		return
	}

	transactions := models.Transactions{}
	err = db.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		AndWhere(dbx.NewExp("`credit`>0")).
		All(&transactions)
	if err != nil {
		c.JSON(500, gin.H{"errors": err.Error()})
		return
	}

	p := profiler.New(transactions, 2)
	res := p.Run()

	c.JSON(200, res)
}
