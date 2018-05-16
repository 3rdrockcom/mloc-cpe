package router

import (
	"net"
	"os"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/controllers"

	"github.com/labstack/echo"
)

type Router struct {
	c *controllers.Controllers
	e *echo.Echo
}

func NewRouter(c *controllers.Controllers) *Router {
	r := &Router{}

	r.e = echo.New()
	r.c = c

	r.appendMiddleware()
	r.appendRoutes()
	r.appendErrorHandler()

	return r
}

func (r *Router) Run() error {
	host := ""
	if os.Getenv("HOST") != "" {
		host = os.Getenv("HOST")
	}

	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	address := net.JoinHostPort(host, port)
	return r.e.Start(address)
}
