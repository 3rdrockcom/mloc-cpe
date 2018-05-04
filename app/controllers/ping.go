package controllers

import (
	"github.com/gin-gonic/gin"
)

func (co *Controllers) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
