package controllers

import (
	"net/http"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"

	"github.com/labstack/echo"
)

func (co *Controllers) Ping(c echo.Context) error {
	c.JSON(http.StatusOK, helpers.H{
		"message": "pong",
	})

	return nil
}
