package routes

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xiaobaiskill/workpool/routes/version"
)

func NewRouter()(router *gin.Engine){
	router = gin.New()

	router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	v1 := router.Group("/v1")
	version.V1RouterMaps(v1)

	return router
}
