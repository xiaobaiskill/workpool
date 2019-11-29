package request

import (
	"encoding/json"
	"fmt"
	"github.com/xiaobaiskill/workpool/pkg/conf"
	"github.com/xiaobaiskill/workpool/pkg/log"
	"github.com/xiaobaiskill/workpool/pool/queue"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 请求的参数
type postValue struct {
	Action string `json:"action"` // 请求网站
	Method string `json:"method"` // 请求类型 get post delete put...
	Data   string `json:"data"`   // 请求的data数据
}

type ResultResp struct {
	*http.Response
	Err error
}

type request struct {
	Proxy  bool
	Result chan ResultResp
	//Url    string `json:"action"`
	//Method string `json:"method"`
	//Query  map[string]interface{}
	*http.Request
}

// 任务结构体
type WorkRequest struct {
	RetryMax int
	Data     request
}

//var add int  // 记录一共运行了多少次

func (w *WorkRequest) Execute() (error) {
	/*defer func() {
		add++
		fmt.Println(add)
	}()*/
	var (
		httpclient    *http.Client
		httpclientmap queue.HttpClientQueue
		err           error
		resp          *http.Response
	)

	if w.Data.Proxy {
		httpclientmap = <-queue.HttpClientQueueChan
		log.Logger.Info("使用的代理ip为：" + httpclientmap.Ip)
		httpclient = httpclientmap.HttpClient
	} else {
		httpclient = &http.Client{Timeout: time.Duration(conf.Conf.Pool.TimeOut) * 1000 * 1000}
	}

	for {

		resp, err = httpclient.Do(w.Data.Request)
		checkOk := w.retryPolicy(resp, err)

		if !checkOk {
			if w.Data.Proxy {
				go func() { queue.HttpClientQueueChan <- httpclientmap }()
			}
			w.Data.Result <- ResultResp{resp, nil}
			break
		}

		if err == nil {
			w.drainBody(resp.Body)
		} else {
			log.Logger.Error(fmt.Sprintf("代理IP：%s,请求失败：%v", httpclientmap.Ip, err))
			if w.Data.Proxy {
				queue.HCMap.Del(httpclientmap.Ip)
				httpclientmap = <-queue.HttpClientQueueChan
				httpclient = httpclientmap.HttpClient
			}
		}

		w.RetryMax--
		if w.RetryMax <= 0 {
			if w.Data.Proxy {
				queue.HCMap.Del(httpclientmap.Ip)
				//go func() { queue.HttpClientQueueChan <- httpclientmap }()
			}
			w.Data.Result <- ResultResp{resp, err}
			break
		}
	}
	return nil
}

func (w *WorkRequest) drainBody(body io.ReadCloser) {
	defer body.Close()

	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, 10))
	if err != nil {
		log.Logger.Error("Error reading response body: " + err.Error())
	}
}

func (w *WorkRequest) retryPolicy(resp *http.Response, err error) (bool) {
	if err != nil {
		return true
	}

	if resp.StatusCode == 0 || resp.StatusCode >= 500 {
		return true
	}
	return false
}

func NewWorkRequest(body []byte) (w *WorkRequest, err error) {
	p := postValue{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Logger.Error("请求数据 有误，无法实现json 转换：" + err.Error())
		return
	}

	r := request{}
	method := "GET"
	switch strings.ToUpper(p.Method) {
	case "POST":
		fallthrough
	case "PUT":
		fallthrough
	case "DELETE":
		method = strings.ToUpper(p.Method)
	}

	r.Request, err = http.NewRequest(method, p.Action, nil)

	if err != nil {
		log.Logger.Error(fmt.Sprintf("new request err: %v", err))
		return
	}

	m := make(map[string]interface{})
	if len(p.Data) > 0 {
		err = json.Unmarshal([]byte(p.Data), &m)
		if err != nil {
			log.Logger.Error("请求的参数中 query 数据json 解析失败：" + err.Error())
			return
		}
	}

	w = &WorkRequest{conf.Conf.Pool.RetryMax, r}
	return

}
