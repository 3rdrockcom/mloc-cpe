package router

import (
	"net/http"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/controllers"

	"github.com/labstack/echo"
)

func (r *Router) appendErrorHandler() {
	r.e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := err.Error()
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		}

		controllers.SendErrorResponse(c, code, message)
		c.Logger().Error(err)
	}
}
