package proxypool

import "net/http"

type HTTPClientMap struct {
	Ip string
	*http.Client
}

type Proxy interface {
	Pop() (HTTPClientMap, bool) // 取数据
	Push(HTTPClientMap)         // 推送数据
	Del(string)                 // 删除数据
	Len() int                   // 代理池数量
	AddMetric(*Metrics)          // 给每一个代理池添加监控
}
