package auth

import (
	API "github.com/epointpayment/mloc-cpe/app/services/api"
	"github.com/juju/errors"

	"github.com/labstack/echo"
)

// BasicValidator is a validator used for basic auth middleware
func BasicValidator(username, password string, c echo.Context) (isValid bool, err error) {
	// Initialize API service
	sa := API.New()

	// Check is user is authorized
	isValid, err = sa.DoAuth(username, password)
	if err != nil {
		err = errors.Trace(err)
	}
	return
}
