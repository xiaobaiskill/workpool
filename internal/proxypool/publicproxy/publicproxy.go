package publicproxy

import (
	. "github.com/xiaobaiskill/workpool/internal/proxypool"
	"github.com/xiaobaiskill/workpool/pkg/log"
	"github.com/xiaobaiskill/workpool/pkg/models"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type publicProxy struct {
	name string
	minSize          int
	httpClients      map[string]HTTPClientMap
	httpClientsQueue chan HTTPClientMap
	getIpChan        chan bool // 是否获取ip
	sync.Mutex
	m *Metrics
}

func (p *publicProxy) init() {
	p.httpClients = make(map[string]HTTPClientMap)
	p.httpClientsQueue = make(chan HTTPClientMap, p.minSize * 4)
	p.getIpChan = make(chan bool)

	p.getPublicProxy()

	go func() {
		tc := time.NewTicker(5 * time.Second)
		defer tc.Stop()
		for {
			select {
			case <-tc.C:
				if p.Len() > p.minSize*3 {
					continue
				}
				p.getRandomPublicProxy()
			case <-p.getIpChan:
				if len(p.httpClients) > p.minSize {
					continue
				}
				p.getRandomPublicProxy()
			}
		}
	}()
}

func (p *publicProxy) Pop() (HTTPClientMap, bool) {
	if p.Len() <= 0 {
		return HTTPClientMap{}, false
	}
	httpclientip := <-p.httpClientsQueue
	return httpclientip, true
}

func (p *publicProxy) Push(httpclientip HTTPClientMap) {
	p.Lock()
	go func() {
		p.httpClientsQueue <- httpclientip
	}()
	p.Unlock()
}

func (p *publicProxy) Del(ip string) {
	p.Lock()
	delete(p.httpClients, ip)
	p.m.ProxypoolnumDec(p.name)
	p.Unlock()

	if len(p.httpClients) < p.minSize {
		go func() {
			p.getIpChan <- true
		}()
	}
}

func (p *publicProxy) Len() int {
	return len(p.httpClients)
}
func (p *publicProxy) AddMetric(mc *Metrics){
	p.m = mc
	p.init()
}

// 将代理ip 放入至 池中
func (p *publicProxy) add(ip *models.IP) {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.httpClients[ip.Data]; ok {
		return
	}

	hCIp := HTTPClientMap{ip.Data, p.createHttpClient(ip)}
	p.httpClients[ip.Data] = hCIp
	p.m.ProxypoolnumInc(p.name)
	go func() {
		p.httpClientsQueue <- hCIp
	}()

}

// 获取免费代理ip
func (p *publicProxy) getPublicProxy() {
	ips, err := models.Conn.GetNumIPWithType(p.minSize)
	if err != nil || len(ips) == 0 {
		// 运气糟糕 没有获取到数据怎么办呀
		//panic("没有获取可用的IP ，程序终止")
		log.Logger.Warn("publicproxy 初次没有获取到代理ip")
		panic("publicproxy 初次没有获取到代理ip")
		return
	}

	for _, ip := range ips {
		p.add(ip)
	}
}

func (p *publicProxy) getRandomPublicProxy() {
	ips, err := models.Conn.GetRandNumIPWithType(p.minSize)
	if err != nil || len(ips) == 0 {
		// 运气糟糕 没有获取到数据怎么办呀
		//panic("没有获取可用的IP ，程序终止")
		log.Logger.Warn("publicproxy 没有获取到代理ip")
		return
	}

	for _, ip := range ips {
		p.add(ip)
	}
}

// 通过免费代理 生成 httpclient
func (p *publicProxy) createHttpClient(ip *models.IP) *http.Client {
	proxy := func(_ *http.Request) (*url.URL, error) {
		http := "http://"

		if ip.Type2 == "https" {
			http = "https://"
		}
		return url.Parse(http + ip.Data)
	}

	transport := &http.Transport{
		Proxy:               proxy,
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 5,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   4 * time.Second,
	}

	return client
}

func NewPublicProxy(minSize int) *publicProxy {
	p := new(publicProxy)
	p.name = "publicproxy"
	p.minSize = minSize

	return p
}
