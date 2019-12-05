package request

import (
	"encoding/json"
	"fmt"
	"github.com/xiaobaiskill/workpool/internal/proxypool"
	"github.com/xiaobaiskill/workpool/pkg/conf"
	"github.com/xiaobaiskill/workpool/pkg/log"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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
	*http.Request
}

// 任务结构体
type WorkRequest struct {
	RetryMax int
	Data     request
}

//var add int  // 记录一共运行了多少次

func (w *WorkRequest) Execute() {
	// 不使用代理 请求
	if !w.Data.Proxy {
		w.notProxyRequest()
	}

	// 使用代理请求
	w.proxyRequest()
}

// 非代理请求
func (w *WorkRequest) notProxyRequest() {
	for {
		httpclient := &http.Client{}
		resp, err := httpclient.Do(w.Data.Request)
		checkOk := w.retryPolicy(resp, err)
		if !checkOk {
			w.Data.Result <- ResultResp{resp, err}
			break
		}
		w.RetryMax--
		if w.RetryMax <= 0 {
			w.Data.Result <- ResultResp{resp, err}
			break
		}
		if err == nil {
			w.drainBody(resp.Body)
		}
	}

}

// 代理请求
func (w *WorkRequest) proxyRequest() {

	var (
		err        error
		resp       *http.Response
		clientMap  proxypool.HTTPClientMap
		httpclient *http.Client
		index      int
		ok         bool
		retryMax   int
	)
	proxypools := proxypool.Newproxypools()
	retryMax = w.RetryMax

	for {
		clientMap, index, ok = proxypools.Pop()
		if !ok {
			w.Data.Result <- ResultResp{resp, fmt.Errorf("not Pop proxyIp")}
			break
		}
		w.RetryMax = retryMax // 重置单个IP重试次数
		log.Logger.Info("使用的代理ip为：" + clientMap.Ip)
		httpclient = clientMap.Client

		for {
			resp, err = httpclient.Do(w.Data.Request)
			checkOk := w.retryPolicy(resp, err)

			if !checkOk {
				go func() { proxypools.Push(index, clientMap) }()
				w.Data.Result <- ResultResp{resp, err}
				return
			}

			if err == nil {
				w.drainBody(resp.Body)
			} else {
				log.Logger.Error(fmt.Sprintf("代理IP：%s,请求失败：%v", clientMap.Ip, err))
				go func() { proxypools.Del(index, clientMap) }()
				break
			}

			w.RetryMax--
			// 单个代理ip的使用次数用完了 ，但是还是没有获取到数据，则结束这个ip ,用下一个
			if w.RetryMax <= 0 {
				go func() { proxypools.Push(index, clientMap) }()
				break
			}
		}
	}

	w.Data.Result <- ResultResp{resp, err}
}

func (w *WorkRequest) drainBody(body io.ReadCloser) {
	defer body.Close()

	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, 10))
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Error reading response body: %v", err))
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
	log.Logger.Info("========== 有代理请求进来了"+ string(body))

	p := postValue{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("请求数据 有误，无法实现json 转换：", err))
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
			log.Logger.Error(fmt.Sprintf("请求的参数中 query 数据json 解析失败：%v", err))
			return
		}
	}

	w = &WorkRequest{conf.Conf.Pool.RetryMax, r}
	log.Logger.Info("=========work生成==========")
	return

}
