package feiyi

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

type Feiyi struct {
}

func New() *Feiyi {
	return &Feiyi{}
}

//feiyi get ip from feiyiproxy.com
func (f *Feiyi) Execute() (result []*models.IP) {
	clog.Info("FEIYI] start test")
	pollURL := "http://www.feiyiproxy.com/?page_id=1457"
	doc, err := htmlquery.LoadURL(pollURL)
	if err != nil {
		clog.Info("FEIYI] LoadURL error")
		clog.Warn(err.Error())
	}
	trNode, err := htmlquery.Find(doc, "//div[@class='et_pb_code et_pb_module  et_pb_code_1']//div//table//tbody//tr")
	clog.Info("[FEIYI] start up")
	if err != nil {
		clog.Info("FEIYI] parse pollUrl error")
		clog.Warn(err.Error())
	}
	//debug begin
	clog.Info("[FEIYI] len(trNode) = %d ", len(trNode))
	for i := 1; i < len(trNode); i++ {
		tdNode, _ := htmlquery.Find(trNode[i], "//td")
		if len(tdNode) < 7 {
			continue
		}
		ip := htmlquery.InnerText(tdNode[0])
		port := htmlquery.InnerText(tdNode[1])
		Type := htmlquery.InnerText(tdNode[3])
		speed := htmlquery.InnerText(tdNode[6])

		IP := models.NewIP()
		IP.Data = ip + ":" + port

		if Type == "HTTPS" {
			IP.Type1 = "https"
			IP.Type2 = ""

		} else if Type == "HTTP" {
			IP.Type1 = "http"
		}
		IP.Speed = extractSpeed(speed)

		clog.Info("[FEIYI] ip.Data = %s,ip.Type = %s,%s ip.Speed = %d", IP.Data, IP.Type1, IP.Type2, IP.Speed)

		result = append(result, IP)
	}

	clog.Info("FEIYI done.")
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
