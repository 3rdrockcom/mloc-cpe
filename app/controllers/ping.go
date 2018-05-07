package controllers

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"

	"github.com/labstack/echo"
)

func (co *Controllers) Ping(c echo.Context) error {
	c.JSON(200, helpers.H{
		"message": "pong",
	})

	return nil
}
