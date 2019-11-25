package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaobaiskill/workpool/pool/queue"
	"net/http"
)

func Ping(c *gin.Context) {
	c.String(200, "pong")
}

func CatProxyIp(c *gin.Context){
	c.JSON(http.StatusOK,gin.H{"useing proxy ip total":queue.HCMap.Len()})
}
