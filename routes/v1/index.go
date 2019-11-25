package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaobaiskill/workpool/controllers/index"
	"github.com/xiaobaiskill/workpool/controllers/ping"
)

func RouterIndexMap(router *gin.RouterGroup){
	router.GET("/ping",ping.Ping)
	router.GET("/catproxyip", ping.CatProxyIp)

	router.POST("/index",index.Index)
}
