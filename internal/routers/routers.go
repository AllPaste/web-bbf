package routers

import (
	"net/http"

	"github.com/AllPaste/web-bbf/config"
	"github.com/AllPaste/web-bbf/pkg/middleware/dump"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(config.Cfg.Server.Mode)

	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	r.Use(dump.Dump())

	// apiv1 := r.Group("/api/v1")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r
}
