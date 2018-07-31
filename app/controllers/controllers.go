package controllers

import (
	"github.com/epointpayment/mloc-cpe/app/codes"
	"github.com/epointpayment/mloc-cpe/app/database"
	"github.com/epointpayment/mloc-cpe/app/helpers"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/juju/errors"
	"github.com/labstack/echo"
)

// DB is the database handler
var db *dbx.DB

// Controllers manages the controllers used in the application
type Controllers struct{}

// NewControllers creates an instance of the service
func NewControllers(database *database.Database) *Controllers {
	db = database.GetInstance()
	c := &Controllers{}
	return c
}

// SendResponse sends a response to requestor
func SendResponse(c echo.Context, code int, i interface{}) (err error) {
	err = c.JSON(code, i)
	if err != nil {
		err = errors.Trace(err)
	}

	return
}

// SendOKResponse sends a StatusOK (200) response to requestor
func SendOKResponse(c echo.Context, res codes.Code) (err error) {
	err = c.JSON(res.StatusCode, helpers.H{
		"status":        true,
		"response_code": res.StatusCode,
		// "code":          res.Name,
		"message": res.Message,
	})
	if err != nil {
		err = errors.Trace(err)
	}

	return
}

// SendErrorResponse sends an error response to requestor
func SendErrorResponse(c echo.Context, res codes.Code) (err error) {
	err = c.JSON(res.StatusCode, helpers.H{
		"status":        false,
		"response_code": res.StatusCode,
		// "code":          res.Name,
		"error": res.Message,
	})
	if err != nil {
		err = errors.Trace(err)
	}

	return
}
