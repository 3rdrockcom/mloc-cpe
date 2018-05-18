package router

import (
	"net/http"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/controllers"
	API "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/api"
	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"

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

		switch message {
		case API.ErrInvalidAPIKey.Error():
			code = http.StatusForbidden
		case Customer.ErrInvalidUniqueCustomerID.Error():
			code = http.StatusForbidden
		case Customer.ErrCustomerNotFound.Error():
			code = http.StatusNotFound
		}

		controllers.SendErrorResponse(c, code, message)
		c.Logger().Error(err)
	}
}
