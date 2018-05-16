package controllers

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"
	"github.com/labstack/echo"
)

type Controllers struct{}

func NewControllers() *Controllers {
	c := &Controllers{}
	return c
}

func SendResponse(c echo.Context, code int, i interface{}) error {
	return c.JSON(code, i)
}

func SendErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, helpers.H{
		"status":        false,
		"response_code": code,
		"error":         message,
	})
}
