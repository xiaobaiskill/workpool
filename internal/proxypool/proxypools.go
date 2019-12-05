package proxypool

type proxypools struct {
	proxypools []proxypoolRetrynum
}

func (r *proxypools) Pop() (clientMap HTTPClientMap, index int,ok bool) {
	var v proxypoolRetrynum
	for index, v = range r.proxypools {
		switch {
		// -1 随便取
		case v.retrynum == -1:
			clientMap, ok = v.proxy.Pop()

			if ok {
				// 存在返回
				return
			} else {
				// 不存在跳过
				continue
			}
		// >0 有次数限制
		case v.retrynum > 0:
			clientMap, ok = v.proxy.Pop()
			if ok {
				// 存在返回  次数减1
				//v.retrynum--
				r.proxypools[index].retrynum--
				return
			} else {
				// 不存在 跳过
				continue
			}
		// 0 这个代理池不能取了
		default:
			continue
		}
	}

	return
}

func (r *proxypools) Push(index int, clientMap HTTPClientMap) {
	proxypoolRetrynums[index].proxy.Push(clientMap)
}

func (r *proxypools) Del(index int, clientMap HTTPClientMap) {
	proxypoolRetrynums[index].proxy.Del(clientMap.Ip)
}

func Newproxypools() *proxypools {
	r := new(proxypools)
	r.proxypools = make([]proxypoolRetrynum,len(proxypoolRetrynums))
	copy(r.proxypools, proxypoolRetrynums)
	return r
}
