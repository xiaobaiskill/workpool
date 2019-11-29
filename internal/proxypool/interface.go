package proxypool

import "net/http"

type HTTPClientMap struct{
	ip string
	*http.Client
}

type Proxy interface {
	Init()     // 初始化
	Pop()      // 取数据
	Push()     // 推送数据
	Del()      // 删除数据
	Add()      // 添加数据
}


