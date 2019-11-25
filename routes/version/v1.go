package version

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/xiaobaiskill/workpool/routes/v1"
)

func V1RouterMaps(router *gin.RouterGroup){
	v1.RouterIndexMap(router)
}
