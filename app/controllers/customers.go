package controllers

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"

	"github.com/gin-gonic/gin"
)

func (co Controllers) PostAddCustomer(c *gin.Context) {
	var err error
	db := co.DB.GetInstance()

	customer := models.Customer{}

	if err = c.ShouldBindJSON(&customer); err != nil {
		c.JSON(400, gin.H{"errors": err.Error()})
		return
	}

	if err = customer.Validate(); err != nil {
		c.JSON(400, gin.H{"errors": err.Error()})
		return
	}

	err = db.Model(&customer).Insert()
	if err != nil {
		c.JSON(500, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(200, customer)
}
