package queue

import (
	"github.com/xiaobaiskill/workpool/pkg/conf"
	"github.com/xiaobaiskill/workpool/pkg/models"
	"github.com/xiaobaiskill/workpool/pool/worker"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	HttpClientQueueChan chan HttpClientQueue
	GetIpChan           = make(chan bool)
	HCMap               HttpClientMap
)

func InitQueue() {
	HCMap = HttpClientMap{}
	HCMap.httpClients = make(map[string]*http.Client)
	HttpClientQueueChan = make(chan HttpClientQueue, conf.Conf.Pool.ProxyIpSize * 3)
}

type HttpClientQueue struct {
	Ip         string
	HttpClient *http.Client
}

type HttpClientMap struct {
	httpClients map[string]*http.Client
	lock        sync.Mutex
}

func (h *HttpClientMap) Len() int {
	return len(h.httpClients)
}

func (h *HttpClientMap) Add(ip *models.IP) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if _, ok := h.httpClients[ip.Data]; ok {
		return
	}

	hCIp := createHttpClient(ip)
	h.httpClients[ip.Data] = hCIp
	go func() {
		HttpClientQueueChan <- HttpClientQueue{ip.Data, hCIp}
	}()
}

func (h *HttpClientMap) Del(ip string) {
	h.lock.Lock()
	if _, ok := h.httpClients[ip]; ok {
		delete(h.httpClients, ip)
	}
	h.lock.Unlock()

	// 保证work 的数量小于代理ip的数量即可
	if len(worker.Works) > h.Len() {
		go func() {GetIpChan<-true}()
	}
}

// 通过IP 生成 httpclient
func createHttpClient(ip *models.IP) *http.Client {
	proxy := func(_ *http.Request) (*url.URL, error) {
		http := "http://"

		if ip.Type2 == "https" {
			http = "https://"
		}
		return url.Parse(http + ip.Data)
	}

	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(conf.Conf.Pool.TimeOut) * 1000 * 1000, //  1 second == 1000 000 000 time.Duration
	}

	return client
}
