package data5u

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/ruoklive/proxypool/pkg/models"
	"github.com/ruoklive/proxypool/pkg/register"
	"log"
	"strconv"
)

func init() {
	register.Add(func() register.IPGetter {
		return New()
	})
}

type Data5u struct {
}

func New() *Data5u {
	return &Data5u{}
}

//Data5u is not work now
// Data5u get ip from data5u.com
func (d *Data5u) Execute() (result []*models.IP) {
	pollURL := "http://www.data5u.com/free/index.shtml"
	resp, _, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return
	}
	if resp.StatusCode != 200 {
		log.Println(errs)
		return
	}
	fmt.Println(resp.Body)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}

	doc.Find("body > div.wlist > ul > li:nth-child(2) > ul").Each(func(i int, s *goquery.Selection) {
		node := strconv.Itoa(i + 1)
		ss := s.Find("ul:nth-child(" + node + ") > span:nth-child(1) > li").Text()
		sss := s.Find("ul:nth-child(" + node + ") > span:nth-child(2) > li").Text()
		ssss := s.Find("ul:nth-child(" + node + ") > span:nth-child(4) > li").Text()
		ip := models.NewIP()
		ip.Data = ss + ":" + sss
		ip.Type1 = ssss
		fmt.Printf("ip.Data = %s, ip.Type = %s", ip.Data, ip.Type1)
		result = append(result, ip)
	})
	log.Println("Data5u done.")
	return
}
