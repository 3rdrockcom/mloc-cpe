package router

import (
	"github.com/epointpayment/mloc-cpe/app/config"
	"github.com/epointpayment/mloc-cpe/app/log"
	"github.com/epointpayment/mloc-cpe/app/router/middleware/auth"
	"github.com/epointpayment/mloc-cpe/app/router/middleware/logger"
	"github.com/epointpayment/mloc-cpe/app/router/middleware/logger/logrus"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// appendMiddleware registers middleware
func (r *Router) appendMiddleware() {
	r.e.Use(middleware.RequestID())
	r.e.Use(middleware.Recover())

	// logger
	r.e.Logger = logrus.New(log.DefaultLogger)
	r.e.Use(logger.LoggerWithConfig(logger.LoggerConfig{
		Logger: log.DefaultLogger,
	}))

	if config.IsDev() {
		r.e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			log.Debug("Request Body:\n" + string(reqBody))
			log.Debug("Response Body:\n" + string(resBody))
		}))
	}
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
	if len(keyLookup) == 0 {
		panic("invalid keyLookup used for middleware CUIDAuth validator")
	}

	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "query:" + keyLookup,
		AuthScheme: "",
		Validator:  auth.CUIDValidator,
	})
}
