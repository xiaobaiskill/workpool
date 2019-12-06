package proxypool

type proxypoolRetrynum struct {
	proxy    Proxy
	retrynum int
}

var proxypoolRetrynums []proxypoolRetrynum

func init() {
	proxypoolRetrynums = make([]proxypoolRetrynum, 0)
}


type Register struct{
	Metrics *Metrics
}

// 注入 不同proxy池，和每个池允许job调用的次数
func (r *Register)Add(proxy Proxy, retrynum int) {
	proxy.AddMetric(r.Metrics)
	proxypoolRetrynums = append(proxypoolRetrynums, proxypoolRetrynum{proxy, retrynum})
}


func NewRegister()*Register{
	r := new(Register)
	r.Metrics = NewMetrics("proxypool")
	return r
}

