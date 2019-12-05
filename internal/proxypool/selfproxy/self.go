package selfproxy

import (
	. "github.com/xiaobaiskill/workpool/internal/proxypool"
	"net/http"
	"strings"
)

type selfproxy struct {
	warnURL string
}

func (s *selfproxy) Pop() (hclient HTTPClientMap, ok bool) {
	client := &http.Client{}
	go func() {
		// 通知啦
		req, err := http.NewRequest("POST", s.warnURL, strings.NewReader(`{"text":"正在使用自身ip 请尽快添加备用ip!!!!!"}`))
		if err != nil {
			return
		}
		req.Header.Set("Content-type", "application/json")
		client.Do(req)
		// 通知结束
	}()

	hclient = HTTPClientMap{"", client}
	ok = true
	return
}

func (s *selfproxy) Push(httpclientip HTTPClientMap) {
	return
}

func (s *selfproxy) Del(ip string) {
	return
}

func (s *selfproxy) Len()int{
	return 1
}

func NewSelf(url string) *selfproxy {
	s := new(selfproxy)
	s.warnURL = url
	return s
}
