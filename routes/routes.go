package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_app/logger"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world!")
	})
	r.GET("/in", func(c *gin.Context) {
		c.String(http.StatusOK, "in hello world!")
	})

	return r
}
