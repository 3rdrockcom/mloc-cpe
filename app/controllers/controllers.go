package controllers

import (
	"net/http"

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

func SendOKResponse(c echo.Context, message string) error {
	code := http.StatusOK
	return c.JSON(code, helpers.H{
		"status":        true,
		"response_code": code,
		"message":       message,
	})
}

func SendErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, helpers.H{
		"status":        false,
		"response_code": code,
		"error":         message,
	})
}
