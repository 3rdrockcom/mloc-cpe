package router

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/router/middleware/auth"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (r *Router) appendMiddleware() {
	r.e.Use(middleware.Gzip())
	r.e.Use(middleware.Logger())
	r.e.Use(middleware.Recover())
}

func (r *Router) mwBasicAuth() echo.MiddlewareFunc {
	return middleware.BasicAuth(auth.BasicValidator)
}

func (r *Router) mwKeyAuth(authType string) echo.MiddlewareFunc {
	validator := auth.DefaultValidator

	switch authType {
	case "login":
		validator = auth.LoginValidator
	case "registration":
		validator = auth.RegistrationValidator
	}

	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:X-API-KEY",
		AuthScheme: "",
		Validator:  validator,
	})
}
