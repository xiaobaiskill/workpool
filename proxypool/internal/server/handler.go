package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-clog/clog"
	"github.com/ruoklive/proxypool/pkg/models"
	"github.com/ruoklive/proxypool/pkg/util"
	"net/http"
)

// random ip
func (s *Server) ProxyRandom(c *gin.Context) {
	ips, err := s.db.GetAllIP()
	if err != nil {
		clog.Warn(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	x := len(ips)
	clog.Info("len(ips) = %d", x)
	if x == 0 {
		c.JSON(http.StatusBadRequest, models.NewIP())
		return
	}
	randomNum := util.RandInt(0, x)
	c.JSON(http.StatusOK, ips[randomNum])

}

func (s *Server) ProxyFind(c *gin.Context) {
	ips, err := s.db.FindIPWithType("https")
	if err != nil {
		clog.Warn(err.Error())
		c.JSON(http.StatusOK, models.NewIP())
		return
	}
	x := len(ips)
	clog.Warn("x = %d", x)
	randomNum := util.RandInt(0, x)
	clog.Info("[proxyFind] random num = %d", randomNum)
	if randomNum == 0 {
		c.JSON(http.StatusOK, models.NewIP())
		return
	}
	c.JSON(http.StatusOK, ips[randomNum])
}

// HealthCheck health check
func (s *Server) HealthCheck(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"msg": "成功",
	})
}

// Monitor
func (s *Server) Monitor(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"waiting": s.workPool.Waiting(),
		"executing": s.workPool.GetStatistics().Executing,
		"total": s.workPool.GetStatistics().Total,
	})
}
