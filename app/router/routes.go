package router

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/controllers"

	"github.com/gin-gonic/gin"
)

func appendRoutes(r *gin.Engine) {
	c := &controllers.Controllers{}

	r.GET("/ping", c.Ping)
}
