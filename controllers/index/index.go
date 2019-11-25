package index

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaobaiskill/workpool/pool/job/request"
	"github.com/xiaobaiskill/workpool/pool/worker"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Index(c *gin.Context){
	if c.Request.Method == "POST" {
		param, err := ioutil.ReadAll(c.Request.Body)
		if len(param) == 0 || err != nil {
			c.String(http.StatusBadRequest,"请携带参数。")
			return
		}

		w,err := request.NewWorkRequest(param)
		if err != nil {
			c.String(http.StatusBadRequest,"请求参数有误，请确认参数格式是否书写正确")
			return
		}

		w.Data.Proxy,_ = strconv.ParseBool(c.Query("proxy"))
		w.Data.Result = make(chan interface{})
		worker.WorkQueue <- w
		meta := <-w.Data.Result
		c.JSON(http.StatusOK,gin.H{"proxy":w.Data.Proxy,"meta":meta})
	} else {
		c.String(http.StatusBadRequest,"method 请求有误")
	}
}

