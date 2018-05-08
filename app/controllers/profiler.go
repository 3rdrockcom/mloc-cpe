package controllers

import (
	"net/http"
	"strconv"
	"time"

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
		c.JSON(http.StatusBadRequest, helpers.H{"errors": err.Error()})
		return nil
	}

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
	q := db.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		AndWhere(dbx.NewExp("`credit`>0")).
		AndWhere(dbx.Between("datetime", startDate, endDate))
	err = q.All(&transactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.H{"errors": err.Error()})
		return nil
	}

	p := profiler.New(transactions, 2)
	res := p.Run()

	c.JSON(http.StatusOK, res)
	return nil
}
