package router

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/epointpayment/mloc-cpe/app/config"
	"github.com/epointpayment/mloc-cpe/app/controllers"
	"github.com/epointpayment/mloc-cpe/app/log"

	"github.com/labstack/echo"
	"github.com/pseidemann/finish"
)

// Router manages the applications routing functions
type Router struct {
	c *controllers.Controllers
	e *echo.Echo
}

// NewRouter creates an instance of the service
func NewRouter(c *controllers.Controllers) *Router {
	r := &Router{}

	r.c = c

	// Initialize router
	r.e = echo.New()
	r.e.Binder = new(CustomBinder)
	r.e.HideBanner = true

	r.appendMiddleware()
	r.appendRoutes()
	r.appendErrorHandler()

	return r
}

func (r *Router) Run() (err error) {
	// Get config information
	host := config.Get().Server.Host
	port := strconv.FormatInt(config.Get().Server.Port, 10)

	// Create an address for the router to use
	r.e.Server.Addr = net.JoinHostPort(host, port)

	// Initialize graceful shutdown service
	gracefully := &finish.Finisher{
		Log:     log.DefaultLogger,
		Timeout: 10 * time.Second,
	}

	// Add a server for gracious shutdown
	gracefully.Add(r.e.Server)

	// Start routing service
	go func() {
		if err := r.e.Server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	gracefully.Wait()
	return
}
