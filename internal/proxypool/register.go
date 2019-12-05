package proxypool

import "time"

type proxypoolRetrynum struct {
	proxy    Proxy
	retrynum int
}

var proxypoolRetrynums []proxypoolRetrynum

func init() {
	proxypoolRetrynums = make([]proxypoolRetrynum, 0)
}


type Register struct{
	Metrics *metrics
}

// 注入 不同proxy池，和每个池允许job调用的次数
func (r *Register)Add(proxy Proxy, retrynum int, proxyName string) {
	r.cronCheckProxyNum(proxyName,proxy)
	proxypoolRetrynums = append(proxypoolRetrynums, proxypoolRetrynum{proxy, retrynum})
}

func (r *Register) cronCheckProxyNum(name string,proxy Proxy){
	go func() {
		tc := time.NewTicker(time.Second)
		for {
			select {
			case <-tc.C:
				r.Metrics.proxypoolnumset(name,float64(proxy.Len()))
			}
		}
	}()
}

func NewRegister()*Register{
	r := new(Register)
	r.Metrics = NewMetrics("proxypool")
	return r
}

