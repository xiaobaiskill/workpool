package index

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xiaobaiskill/workpool/internal/task/request"
	"github.com/xiaobaiskill/workpool/pkg/log"
	"github.com/xiaobaiskill/workpool/pkg/pool"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Index(c *gin.Context){
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
		w.Data.Result = make(chan request.ResultResp,1)
		pool.WorkQueue <- w
		resultresp := <-w.Data.Result
		close(w.Data.Result)
		log.Logger.Debug(fmt.Sprintf("response:%v ,err:%v",resultresp.Response.Body,resultresp.Err))
		if resultresp.Err != nil {
			log.Logger.Info(fmt.Sprintf("请求报错啦：%v",err))
			c.JSON(http.StatusBadRequest,gin.H{"proxy":w.Data.Proxy,"meta":""})
		} else {
			b,err := ioutil.ReadAll(resultresp.Response.Body)
			resultresp.Response.Body.Close()
			if err != nil {
				c.JSON(resultresp.Response.StatusCode,gin.H{"proxy":w.Data.Proxy,"meta":""})
				return
			}
			c.JSON(resultresp.Response.StatusCode,gin.H{"proxy":w.Data.Proxy,"meta":string(b)})
		}
		return
}

