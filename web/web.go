package web

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xgadget-lab/nexttrace/config"
	"github.com/xgadget-lab/nexttrace/ipgeo"
	"github.com/xgadget-lab/nexttrace/trace"
)

var confToken string

func Start() {
	configData, err := config.Read()

	// Initialize Default Config
	if err != nil || configData.DataOrigin == "" {
		if configData, err = config.AutoGenerate(); err != nil {
			log.Fatal(err)
		}
	}

	// Check Token Available
	if configData.APIToken == "" {
		confToken = "NextTrace"
	} else {
		confToken = configData.APIToken
	}

	router := gin.Default()

	router.GET("/trace", func(c *gin.Context) {
		var timeout time.Duration

		token := c.Query("token")
		ip := net.ParseIP(c.Query("ip"))
		m := c.DefaultQuery("method", "icmp")
		dataOrigin := c.DefaultQuery("data", "LeoMoeAPI")

		if token != confToken {
			c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Token错误"})
			return
		}

		if m == "icmp" {
			timeout = 600
		} else {
			timeout = 1100
		}

		var conf = trace.Config{
			DestIP:           ip,
			DestPort:         80,
			MaxHops:          30,
			NumMeasurements:  3,
			ParallelRequests: 18,
			RDns:             true,
			IPGeoSource:      ipgeo.GetSource(dataOrigin),
			Timeout:          timeout * time.Millisecond,
		}

		res, _ := trace.Traceroute(trace.Method(m), conf)

		c.JSON(http.StatusOK, res)
	})

	var port string
	if configData.ListenPort == 0 {
		port = "8080"
	} else {
		port = strconv.Itoa(configData.ListenPort)
	}
	router.Run(":" + port)
}
