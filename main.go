package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/OwO-Network/nexttrace-enhanced/config"
	fastTrace "github.com/OwO-Network/nexttrace-enhanced/fast_trace"
	"github.com/OwO-Network/nexttrace-enhanced/ipgeo"
	"github.com/OwO-Network/nexttrace-enhanced/printer"
	"github.com/OwO-Network/nexttrace-enhanced/reporter"
	"github.com/OwO-Network/nexttrace-enhanced/trace"
	"github.com/OwO-Network/nexttrace-enhanced/tracemap"
	"github.com/OwO-Network/nexttrace-enhanced/util"
	"github.com/OwO-Network/nexttrace-enhanced/web"
	"github.com/OwO-Network/nexttrace-enhanced/wshandle"
	"github.com/syndtr/gocapability/capability"
)

var fSet = flag.NewFlagSet("", flag.ExitOnError)
var webAPI = fSet.Bool("w", false, "Enable Web API Method")
var fastTest = fSet.Bool("f", false, "One-Key Fast Traceroute")
var tcpSYNFlag = fSet.Bool("T", false, "Use TCP SYN for tracerouting (default port is 80)")
var udpPackageFlag = fSet.Bool("U", false, "Use UDP Package for tracerouting (default port is 53 in UDP)")
var port = fSet.Int("p", 80, "Set SYN Traceroute Port")
var manualConfig = fSet.Bool("c", false, "Manual Config [Advanced]")
var numMeasurements = fSet.Int("q", 3, "Set the number of probes per each hop.")
var parallelRequests = fSet.Int("r", 18, "Set ParallelRequests number. It should be 1 when there is a multi-routing.")
var maxHops = fSet.Int("m", 30, "Set the max number of hops (max TTL to be reached).")
var dataOrigin = fSet.String("d", "", "Choose IP Geograph Data Provider [LeoMoeAPI, IP.SB, IPInfo, IPInsight, IPAPI.com]")
var noRdns = fSet.Bool("n", false, "Disable IP Reverse DNS lookup")
var routePath = fSet.Bool("report", false, "Route Path")
var tablePrint = fSet.Bool("table", false, "Output trace results as table")
var ver = fSet.Bool("V", false, "Check Version")
var timeOut = fSet.Int("t", 1000, "Set timeout [Millisecond]")
var fixIPGeoMode = fSet.Bool("fix", false, "Fix IP Geo Mode")
var country = fSet.String("fix-country", "", "Set Country")
var prov = fSet.String("fix-prov", "", "Set Province/Region")
var city = fSet.String("fix-city", "", "Set City/Area")
var beginHop = fSet.Int("b", 1, "Set the begin hop")
var classicPrint = fSet.Bool("classic", false, "Classic Output trace results like BestTrace")
var jsonEnable = fSet.Bool("j", false, "Output with json format")
var ipv4Only = fSet.Bool("4", false, "Only Displays IPv4 addresses")
var ipv6Only = fSet.Bool("6", false, "Only Displays IPv6 addresses")
var maptrace = fSet.Bool("M", false, "No Print Trace Map")
var src_addr = fSet.String("S", "", "Use the following IP address as the source address in outgoing packets")
var src_dev = fSet.String("D", "", "Use the following Network Devices as the source address in outgoing packets")
var dns_ip = fSet.String("dns", "", "Use the following IP address to resolve domain")

func printArgHelp() {
	fmt.Println("\nArgs Error\nUsage : 'nexttrace [option...] HOSTNAME' or 'nexttrace HOSTNAME [option...]'\nOPTIONS: [-VTU] [-d DATAORIGIN.STR ] [ -m TTL ] [ -p PORT ] [ -q PROBES.COUNT ] [ -r PARALLELREQUESTS.COUNT ] [-rdns] [ -table ] -report")
	fSet.PrintDefaults()
	os.Exit(2)
}

func flagApply() string {

	target := ""
	if len(os.Args) < 2 {
		printArgHelp()
	}

	// flag parse
	if !strings.HasPrefix(os.Args[1], "-") {
		target = os.Args[1]
		err := fSet.Parse(os.Args[2:])
		if err != nil {
			return ""
		}
	} else {
		err := fSet.Parse(os.Args[1:])
		if err != nil {
			return ""
		}
		target = fSet.Arg(0)
	}

	if !*jsonEnable {
		printer.Version()
	}

	// Print Version
	if *ver {
		printer.CopyRight()
		os.Exit(0)
	}

	// Advanced Config
	if *manualConfig {
		if err := config.Generate(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if *webAPI {
		web.Start()
	}

	// -f Fast Test
	if *fastTest {
		fastTrace.FastTest(*tcpSYNFlag, *src_dev, *src_addr)
		os.Exit(0)
	}

	if target == "" {
		printArgHelp()
	}
	return target
}

func main() {

	domain := flagApply()

	capabilities_check()

	configData, err := config.Read()

	// Initialize Default Config
	if err != nil || configData.DataOrigin == "" {
		if configData, err = config.AutoGenerate(); err != nil {
			log.Fatal(err)
		}
	}

	// Set Token from Config
	ipgeo.SetToken(configData.Token)

	// Check Whether User has specified IP Geograph Data Provider
	if *dataOrigin == "" {
		// Use Default Data Origin with Config
		*dataOrigin = configData.DataOrigin
	}

	var ip net.IP

	if *tcpSYNFlag || *udpPackageFlag {
		ip = util.DomainLookUp(domain, *dns_ip, true, false, *jsonEnable)
	} else {
		ip = util.DomainLookUp(domain, *dns_ip, *ipv4Only, *ipv6Only, *jsonEnable)
	}

	// if ip.To4() == nil && strings.ToUpper(*dataOrigin) == "LEOMOEAPI" {
	// 	// IPv6 不使用 LeoMoeAPI
	// 	*dataOrigin = "ipinsight"
	// }

	if *src_dev != "" {
		dev, _ := net.InterfaceByName(*src_dev)

		if addrs, err := dev.Addrs(); err == nil {
			for _, addr := range addrs {
				if (addr.(*net.IPNet).IP.To4() == nil) == (ip.To4() == nil) {
					*src_addr = addr.(*net.IPNet).IP.String()
				}
			}
		}
	}

	if strings.ToUpper(*dataOrigin) == "LEOMOEAPI" {
		w := wshandle.New()
		w.Interrupt = make(chan os.Signal, 1)
		signal.Notify(w.Interrupt, os.Interrupt)
		defer func() {
			err := w.Conn.Close()
			if err != nil {
				return
			}
		}()
	}

	if !*jsonEnable {
		printer.PrintTraceRouteNav(ip, domain, *dataOrigin)
	}

	var m trace.Method = ""

	switch {
	case *tcpSYNFlag:
		m = trace.TCPTrace
	case *udpPackageFlag:
		m = trace.UDPTrace
	default:
		m = trace.ICMPTrace
	}

	if !*tcpSYNFlag && *port == 80 {
		*port = 53
	}

	if !*noRdns {
		*noRdns = configData.NoRDNS
	}

	var conf = trace.Config{
		SrcAddr:          *src_addr,
		BeginHop:         *beginHop,
		DestIP:           ip,
		DestPort:         *port,
		MaxHops:          *maxHops,
		NumMeasurements:  *numMeasurements,
		ParallelRequests: *parallelRequests,
		RDns:             !*noRdns,
		IPGeoSource:      ipgeo.GetSource(*dataOrigin),
		Timeout:          time.Duration(*timeOut) * time.Millisecond,
	}

	if !*tablePrint && !configData.TablePrintDefault && !*jsonEnable {
		if *classicPrint {
			conf.RealtimePrinter = printer.ClassicPrinter
		} else {
			conf.RealtimePrinter = printer.RealtimePrinter
		}
	}

	res, err := trace.Traceroute(m, conf)

	if err != nil {
		log.Fatalln(err)
	}

	if *fixIPGeoMode {
		f := ipgeo.FixData{
			Country: *country,
			Prov:    *prov,
			City:    *city,
		}
		var ipSplice []net.Addr
		for _, allHops := range res.Hops {
			for _, ttlHops := range allHops {
				if ttlHops.Address != nil {
					// IP 去重
					tFlag := true
					for _, v := range ipSplice {
						if v.String() == ttlHops.Address.String() {
							tFlag = false
							break
						}
					}
					if tFlag {
						if (ttlHops.Geo.Country != f.Country && f.Country != "") || (ttlHops.Geo.Prov != f.Prov && f.Prov != "") || (ttlHops.Geo.City != f.City && f.City != "") {
							ipSplice = append(ipSplice, ttlHops.Address)
						}
						tFlag = false
					}
				}
			}
		}
		for _, v := range ipSplice {
			ipgeo.UpdateIPGeo(v.String(), f)
		}
	}

	if *jsonEnable {
		r := printer.ParseJson(res)
		fmt.Println(r)
		tracemap.GetMapUrl(r)
		<-time.After(10 * time.Millisecond)
		return
	}

	if *tablePrint || configData.TablePrintDefault {
		printer.TracerouteTablePrinter(res)
	}

	if *routePath || configData.AlwaysRoutePath {
		r := reporter.New(res, ip.String())
		r.Print()
	}

	if !*maptrace {
		r := printer.ParseJson(res)
		tracemap.GetMapUrl(r)
		<-time.After(10 * time.Millisecond)
	}

}

func capabilities_check() {

	// Windows 判断放在前面，防止遇到一些奇奇怪怪的问题
	if runtime.GOOS == "windows" {
		// Running on Windows, skip checking capabilities
		return
	}

	uid := os.Getuid()
	if uid == 0 {
		// Running as root, skip checking capabilities
		return
	}

	/***
	* 检查当前进程是否有两个关键的权限
	==== 看不到我 ====
	* 没办法啦
	* 自己之前承诺的坑补全篇
	* 被迫填坑系列 qwq
	==== 看不到我 ====
	***/

	// NewPid 已经被废弃了，这里改用 NewPid2 方法
	caps, err := capability.NewPid2(0)
	if err != nil {
		// 判断是否为macOS
		if runtime.GOOS == "darwin" {
			// macOS下报错有问题
		} else {
			fmt.Println(err)
		}
		return
	}

	// load 获取全部的 caps 信息
	err = caps.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 判断一下权限有木有
	if caps.Get(capability.EFFECTIVE, capability.CAP_NET_RAW) && caps.Get(capability.EFFECTIVE, capability.CAP_NET_ADMIN) {
		// 有权限啦
		return
	} else {
		// 没权限啦
		fmt.Println("您正在以普通用户权限运行 NextTrace，但 NextTrace 未被赋予监听网络套接字的ICMP消息包、修改IP头信息（TTL）等路由跟踪所需的权限")
		fmt.Println("请使用管理员用户执行 `sudo setcap cap_net_raw,cap_net_admin+eip ${your_nexttrace_path}/nexttrace` 命令，赋予相关权限后再运行~")
		fmt.Println("什么？为什么 ping 普通用户执行不要 root 权限？因为这些工具在管理员安装时就已经被赋予了一些必要的权限，具体请使用 `getcap /usr/bin/ping` 查看")
	}
}
