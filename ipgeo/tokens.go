package ipgeo

import "github.com/OwO-Network/nexttrace-enhanced/config"

type tokenData struct {
	ipinsight string
	ipinfo    string
	ipleo     string
}

var token = tokenData{
	ipinsight: "",
	ipinfo:    "",
	ipleo:     "NextTraceDemo",
}


func SetToken(c config.Token) {
	token.ipleo = c.LeoMoeAPI
	token.ipinfo = c.IPInfo
}
