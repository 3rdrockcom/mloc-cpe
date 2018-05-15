package auth

import (
	API "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/api"
	"github.com/labstack/echo"
)

func BasicValidator(username, password string, c echo.Context) (isValid bool, err error) {
	sa := API.New()

	isValid, err = sa.DoAuth(username, password)

	return
}
