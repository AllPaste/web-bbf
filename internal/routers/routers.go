package routers

import (
	"github.com/AllPaste/web-bbf/config"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(config.Cfg.Server.Mode)

	// apiv1 := r.Group("/api/v1")

	return r
}
