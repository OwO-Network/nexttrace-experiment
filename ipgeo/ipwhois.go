package ipgeo

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func IPWHOIS(ip string) (*IPGeoData, error) {
	url := "https://ipwho.is/" + ip
	client := &http.Client{
		// 2 秒超时
		Timeout: 2 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:100.0) Gecko/20100101 Firefox/100.0")
	content, err := client.Do(req)
	if err != nil {
		log.Println("ipwho.is 请求超时(2s)，请切换其他API使用")
		return nil, err
	}
	body, _ := io.ReadAll(content.Body)
	res := gjson.ParseBytes(body)

	var country string
	var prov string

	if res.Get("country").String() == "Hong Kong" {
		country = "China"
	} else if res.Get("country").String() == "Taiwan" {
		country = "China"
		prov = "Taiwan, " + prov
	} else {
		country = res.Get("country").String()
	}

	return &IPGeoData{
		Asnumber: res.Get("connection").Get("asn").String(),
		Country:  country,
		City:     res.Get("city").String(),
		Prov:     prov + res.Get("region").String(),
		Isp:      res.Get("connection").Get("domain").String(),
	}, nil
}
