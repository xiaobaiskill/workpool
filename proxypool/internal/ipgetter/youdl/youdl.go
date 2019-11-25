package youdl

import (
	"github.com/ruoklive/proxypool/pkg/register"
	"log"
	"strings"

	"github.com/ruoklive/proxypool/pkg/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
)
func init() {
	register.Add(func() register.IPGetter {
		return New()
	})
}
type Youdl struct {

}

func New() *Youdl {
	return &Youdl{}
}

// YDL get ip from youdaili.net
func (y *Youdl) Execute() (result []*models.IP) {
	pollURL := "http://www.youdaili.net/Daili/http/"
	_, body, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return
	}
	do, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Println(err.Error())
		return
	}

	URL, _ := do.Find("body > div.con.PT20 > div.conl > div.lbtc.l > div.chunlist > ul > li:nth-child(1) > p > a").Attr("href")
	_, content, errs := gorequest.New().Get(URL).End()
	if errs != nil {
		log.Println(errs)
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Println(err.Error())
		return
	}
	doc.Find(".content p").Each(func(_ int, s *goquery.Selection) {
		ip := models.NewIP()
		c := strings.Split(s.Text(), "@")
		ip.Data = c[0]
		ip.Type1 = strings.ToLower(strings.Split(c[1], "#")[0])
		result = append(result, ip)
	})
	log.Println("YDL done.")
	return
}
