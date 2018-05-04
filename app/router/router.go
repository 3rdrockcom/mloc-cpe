package router

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/controllers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	c *controllers.Controllers
	e *gin.Engine
}

func NewRouter(c *controllers.Controllers) *Router {
	r := &Router{}

	r.e = gin.Default()
	r.c = c
	r.appendRoutes()

	return r
}

func (r *Router) Run() error {
	return r.e.Run()
}
