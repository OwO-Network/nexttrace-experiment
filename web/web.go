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

func CrosHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "false")
		context.Set("content-type", "application/json")

		//处理请求
		context.Next()
	}
}

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

	router.LoadHTMLFiles("web/templates/index.tmpl")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "NextTrace 路由跟踪测试",
			"token": confToken,
		})
	})

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
			NumMeasurements:  1,
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