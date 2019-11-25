package ip66

import (
	"github.com/go-clog/clog"
	"github.com/ruoklive/proxypool/pkg/models"
	"github.com/ruoklive/proxypool/pkg/register"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func init() {
	register.Add(func() register.IPGetter {
		return New()
	})
}

type Ip66 struct {
}

func New() *Ip66 {
	return &Ip66{}
}

// IP66 get ip from 66ip.cn
func (i *Ip66) Execute() (result []*models.IP) {
	var ExprIP = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)

	pollURL := "http://www.66ip.cn/mo.php?tqsl=100"
	resp, err := http.Get(pollURL)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println(err)
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
		clog.Info("[IP66] ip = %s, type = %s", ip.Data, ip.Type1)
		result = append(result, ip)
	}

	clog.Info("IP66 done.")
	return
}
