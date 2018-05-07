package controllers

import (
	"strconv"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/profiler"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/echo"
)

func (co Controllers) GetCustomerProfile(c echo.Context) error {
	var err error
	db := co.DB.GetInstance()

	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		c.JSON(400, helpers.H{"errors": err.Error()})
		return nil
	}

	transactions := models.Transactions{}
	err = db.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		AndWhere(dbx.NewExp("`credit`>0")).
		All(&transactions)
	if err != nil {
		c.JSON(500, helpers.H{"errors": err.Error()})
		return nil
	}

	p := profiler.New(transactions, 2)
	res := p.Run()

	c.JSON(200, res)
	return nil
}
