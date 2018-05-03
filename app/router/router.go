package router

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	e *gin.Engine
}

func NewRouter() *Router {
	r := &Router{}

	r.e = gin.Default()
	appendRoutes(r.e)

	return r
}

func (r *Router) Run() error {
	return r.e.Run()
}
