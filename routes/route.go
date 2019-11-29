package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaobaiskill/workpool/routes/version"
)

func NewRouter(router *gin.Engine){
	v1 := router.Group("/v1")
	version.V1RouterMaps(v1)

}
