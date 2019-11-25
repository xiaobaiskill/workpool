package kuaidl

import (
	"github.com/Aiicy/htmlquery"
	"github.com/go-clog/clog"
	"github.com/ruoklive/proxypool/pkg/models"
	"github.com/ruoklive/proxypool/pkg/register"
	"regexp"
	"strconv"
)

func init() {
	register.Add(func() register.IPGetter {
		return New()
	})
}

type Kuaidl struct {
}

func New() *Kuaidl {
	return &Kuaidl{}
}

// KDL get ip from kuaidaili.com
func (k *Kuaidl) Execute() (result []*models.IP) {
	pollURL := "http://www.kuaidaili.com/free/inha/"
	doc, _ := htmlquery.LoadURL(pollURL)
	trNode, err := htmlquery.Find(doc, "//table[@class='table table-bordered table-striped']//tbody//tr")
	if err != nil {
		clog.Warn(err.Error())
	}
	for i := 0; i < len(trNode); i++ {
		tdNode, _ := htmlquery.Find(trNode[i], "//td")
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		Type := htmlquery.InnerText(tdNode[3])
		speed := htmlquery.InnerText(tdNode[5])

		IP := models.NewIP()
		IP.Data = ip + ":" + port
		if Type == "HTTPS" {
			IP.Type1 = ""
			IP.Type2 = "https"
		} else if Type == "HTTP" {
			IP.Type1 = "http"
		}
		IP.Speed = extractSpeed(speed)
		clog.Info("IP:%s", IP)
		result = append(result, IP)
	}

	clog.Info("[kuaidaili] done")
	return
}

func extractSpeed(oritext string) int64 {
	reg := regexp.MustCompile(`\d+?\.?\d*`)
	temp := reg.FindStringSubmatch(oritext)

	if len(temp) >= 1 && temp[0] != "" {
		speedFloat, _ := strconv.ParseFloat(temp[0], 64)
		speed := int64(speedFloat * 1000)
		return speed
	}
	return -1
}
