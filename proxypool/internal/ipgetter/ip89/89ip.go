package ip89

import (
	"github.com/ruoklive/proxypool/pkg/register"
	"io/ioutil"
	"net/http"
	//"fmt"
	"github.com/go-clog/clog"

	"github.com/ruoklive/proxypool/pkg/models"
	"regexp"
	"strings"
)

func init() {
	register.Add(func() register.IPGetter {
		return New()
	})
}

type Ip89 struct {
}

func New() *Ip89 {
	return &Ip89{}
}

//IP89 get ip from www.89ip.cn
func (i *Ip89) Execute() (result []*models.IP) {
	clog.Info("89IP] start test")
	var ExprIP = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)
	pollURL := "http://www.89ip.cn/tqdl.html?api=1&num=100&port=&address=%E7%BE%8E%E5%9B%BD&isp="

	resp, err := http.Get(pollURL)
	if err != nil {
		clog.Warn(err.Error())
		return
	}

	if resp.StatusCode != 200 {
		clog.Warn(err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyIPs := string(body)
	ips := ExprIP.FindAllString(bodyIPs, 100)

	for index := 0; index < len(ips); index++ {
		ip := models.NewIP()
		ip.Data = strings.TrimSpace(ips[index])
		ip.Type1 = "http"
		clog.Info("[89IP] ip = %s, type = %s", ip.Data, ip.Type1)
		result = append(result, ip)
	}

	clog.Info("89IP done.")
	return
}
