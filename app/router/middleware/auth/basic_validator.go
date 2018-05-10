package auth

import "github.com/labstack/echo"

func BasicValidator(username, password string, c echo.Context) (bool, error) {
	if username == "joe" && password == "secret" {
		return true, nil
	}
	return false, nil
}
