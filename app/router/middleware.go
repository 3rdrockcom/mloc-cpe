package router

import (
	"strings"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/router/middleware/auth"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// appendMiddleware registers middleware
func (r *Router) appendMiddleware() {
	r.e.Use(middleware.Gzip())
	r.e.Use(middleware.Logger())
	r.e.Use(middleware.Recover())
}

// mwBasicAuth handles the basic authentication for a specific route
func (r *Router) mwBasicAuth() echo.MiddlewareFunc {
	return middleware.BasicAuth(auth.BasicValidator)
}

// mwKeyAuth handles the key authentication for a specific route
func (r *Router) mwKeyAuth(authKeyType string) echo.MiddlewareFunc {
	var validator middleware.KeyAuthValidator

	// Assign appropriate key auth validator
	switch authKeyType {
	case "login":
		validator = auth.LoginValidator
	case "registration":
		validator = auth.RegistrationValidator
	case "customer":
		validator = auth.CustomerValidator
	default:
		panic("unknown authKeyType used for middleware KeyAuth validator")
	}

	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:X-API-KEY",
		AuthScheme: "",
		Validator:  validator,
	})
}

// mwCUIDAuth handles the customer unique ID (CUID) authentication for a specific route
func (r *Router) mwCUIDAuth(keyLookup string) echo.MiddlewareFunc {
	// Check if keyLookup value looks valid
	parts := strings.Split(keyLookup, ":")
	if !(len(parts) == 2 && ((parts[0] == "query" || parts[0] == "form") && len(parts[1]) > 0)) {
		panic("invalid keyLookup used for middleware CUIDAuth validator")
	}

	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  keyLookup,
		AuthScheme: "",
		Validator:  auth.CUIDValidator,
	})
}
