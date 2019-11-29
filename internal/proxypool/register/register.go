package register

import (
	. "github.com/xiaobaiskill/workpool/internal/proxypool"
)

type proxypoolnum struct {
	proxy Proxy
	retrynum int
}

var proxypoolnums []proxypoolnum

func Add(proxy Proxy,retrynum int){
	proxypoolnums = append(proxypoolnums,proxypoolnum{proxy,retrynum})
}

func Pop(){
	for _,v := range proxypoolnums{
		switch  {
		case v.retrynum == -1 :
			v.proxy.Pop()
		case v.retrynum > 0 :
			v.proxy.Pop()
			v.retrynum--
		default:
			// case v.retrynum == 0
			continue
		}
	}
}


func Push(id int){
	proxypoolnums[id].proxy.Push()
}

func Del(id int){
	proxypoolnums[id].proxy.Del()
}