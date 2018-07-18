package router

import (
	"net/http"

	"github.com/epointpayment/mloc-cpe/app/controllers"
	API "github.com/epointpayment/mloc-cpe/app/services/api"
	Customer "github.com/epointpayment/mloc-cpe/app/services/customer"

	"github.com/labstack/echo"
)

// appendErrorHandler handles errors for the router
func (r *Router) appendErrorHandler() {
	r.e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := err.Error()
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		}

		// Override status code based on error responses
		switch message {
		// API Service
		case API.ErrInvalidAPIKey.Error():
			code = http.StatusForbidden

		// Customer Service
		case Customer.ErrInvalidUniqueCustomerID.Error():
			code = http.StatusForbidden
		case Customer.ErrCustomerNotFound.Error():
			code = http.StatusNotFound
		case Customer.ErrInvalidData.Error():
			code = http.StatusBadRequest

		// Auth Middleware
		case "missing key in request header":
			message = API.ErrMissingAPIKey.Error()
		case "missing key in the query string":
			message = Customer.ErrMissingUniqueCustomerID.Error()
		case "missing key in the form":
			message = Customer.ErrMissingUniqueCustomerID.Error()

		// Unknown error
		default:
			if _, ok := err.(*echo.HTTPError); !ok {
				message = "Internal Error"
			}
		}

		// Send error in a specific format
		controllers.SendErrorResponse(c, code, message)

		// Log errors
		c.Logger().Error(err)
	}
}
