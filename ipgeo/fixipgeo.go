package ipgeo

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xgadget-lab/nexttrace/config"
)

type FixData struct {
	Country string
	Prov    string
	City    string
}

func UpdateIPGeo(ip string, fixData FixData) {
	configData, err := config.Read()

	// Initialize Default Config
	if err != nil || configData.DataOrigin == "" {
		if configData, err = config.AutoGenerate(); err != nil {
			log.Fatal(err)
		}
	}
	url := fmt.Sprintf("https://api.leo.moe/ip/update.php?ip=%s&country=%s&prov=%s&city=%s&ut=%s", ip, fixData.Country, fixData.Prov, fixData.City, configData.LeoMoeUpdateKey)
	client := &http.Client{
		// 2 秒超时
		Timeout: 2 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "NextTrace Enhanced 0.1.17")
	_, err = client.Do(req)
	if err != nil {
		log.Println("Update 超时")
		return
	}
	fmt.Println("修复IP: " + ip)
}
