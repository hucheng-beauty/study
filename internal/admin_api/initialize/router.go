package initialize

import (
	"study/internal/admin_api/middlewares"
	"study/internal/admin_api/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	Router.Use(middlewares.CORS, middlewares.LimitHandler(maxConn), gin.Logger(), gin.Recovery())

	Router.POST("/", router.Handler)

	return Router
}
