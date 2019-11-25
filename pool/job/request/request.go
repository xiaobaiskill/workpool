package request

import (
	"encoding/json"
	"github.com/xiaobaiskill/workpool/pkg/conf"
	"github.com/xiaobaiskill/workpool/pkg/log"
	"github.com/xiaobaiskill/workpool/pool/queue"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 请求的参数
type PostValue struct {
	Action string `json:"action"` // 请求网站
	Method string `json:"method"` // 请求类型 get post delete put...
	Data   string `json:"data"`   // 请求的data数据
}

type Request struct {
	Proxy  bool
	Url    string `json:"action"`
	Method string `json:"method"`
	Query  map[string]interface{}
	Result chan interface{}
}


// 任务结构体
type WorkRequest struct {
	Data Request
}

func (w *WorkRequest) Execute() (err error ){
	var (
		httpclient *http.Client
		httpclientmap queue.HttpClientQueue
		body []byte
		)
	defer func() {
		w.Data.Result <- string(body)
	}()

	if w.Data.Proxy {
		httpclientmap = <-queue.HttpClientQueueChan
		log.Logger.Info("使用的代理ip为：" + httpclientmap.Ip)
		httpclient = httpclientmap.HttpClient
	} else {
		httpclient = &http.Client{Timeout: time.Duration(conf.Conf.Pool.TimeOut) * 1000 * 1000}
	}


	req, err := http.NewRequest(w.Data.Method, w.Data.Url, nil)
	if err != nil {
		return
	}
	if len(w.Data.Query) != 0 {
		q := req.URL.Query()
		for k, v := range w.Data.Query {
			q.Add(k, v.(string))
		}
		req.URL.RawQuery = q.Encode()
	}

	var resp *http.Response
	resp, err = httpclient.Do(req)
	if err != nil {
		log.Logger.Error("请求出错："+ err.Error())
		if w.Data.Proxy{
			queue.HCMap.Del(httpclientmap.Ip)
		}

		return
	}
	defer resp.Body.Close()
	if w.Data.Proxy {
		go func() { queue.HttpClientQueueChan <- httpclientmap }()
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Error("BODY 获取数据出错："+ err.Error())
	}

	return
}

func NewWorkRequest(body []byte) (w *WorkRequest, err error) {
	p := PostValue{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Logger.Error("请求数据 有误，无法实现json 转换：" + err.Error())
		return
	}

	r := Request{}
	switch strings.ToUpper(p.Method) {
	case "POST":
		fallthrough
	case "PUT":
		fallthrough
	case "DELETE":
		r.Method = strings.ToUpper(p.Method)
	default:
		r.Method = "GET"
	}

	r.Url = p.Action
	m := make(map[string]interface{})
	if len(p.Data) > 0{
		err = json.Unmarshal([]byte(p.Data), &m)
		if err != nil {
			log.Logger.Error("请求的参数中 query 数据json 解析失败：" + err.Error())
			return
		}
	}

	r.Query = m
	w = &WorkRequest{r}
	return

}
