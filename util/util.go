package util

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/fatih/color"
)

// get the local ip and port based on our destination ip
func LocalIPPort(dstip net.IP) (net.IP, int) {
	serverAddr, err := net.ResolveUDPAddr("udp", dstip.String()+":12345")
	if err != nil {
		log.Fatal(err)
	}

	// We don't actually connect to anything, but we can determine
	// based on our destination ip what source ip we should use.
	if con, err := net.DialUDP("udp", nil, serverAddr); err == nil {
		defer con.Close()
		if udpaddr, ok := con.LocalAddr().(*net.UDPAddr); ok {
			return udpaddr.IP, udpaddr.Port
		}
	}
	return nil, -1
}

func DomainLookUp(host string, customDNS string, ipv4Only bool, ipv6Only bool, auto bool) net.IP {
	// 手动构造 Resolver 以定制化 DNS 服务器 IP 等参数
	r := &net.Resolver{
		// 尽管编译器已经禁用 Cgo，这里以防万一，保证无论何种编译环境下都能优先使用 Pure-Go，构造详见 lookup.go 源码
		PreferGo: true,
	}

	if customDNS != "" {
		r.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}
			// 见文档 - Dial uses context.Background internally; to specify the context, use DialContext.
			return d.DialContext(ctx, "udp", customDNS+":53")
		}
	}

	ips, err := r.LookupHost(context.Background(), host)
	if err != nil {
		fmt.Println("Domain " + host + " Lookup Fail.")
		os.Exit(1)
	}

	var ipSlice = []net.IP{}
	var ipv6Flag = false

	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ipv4Only {
			// 仅返回ipv4的ip
			if ip.To4() != nil {
				ipSlice = append(ipSlice, ip)
			} else {
				ipv6Flag = true
			}
		} else if ipv6Only {
			if ip.To4() == nil {
				ipSlice = append(ipSlice, ip)
			}
		} else {
			ipSlice = append(ipSlice, ip)
		}
	}

	if ipv6Flag {
		if !auto {
			// fmt.Println("[Info] IPv6 TCP/UDP Traceroute is not supported right now.")
		}

		if len(ipSlice) == 0 {
			os.Exit(0)
		}
	}

	if len(ipSlice) == 1 || auto {
		return ipSlice[0]
	} else {
		fmt.Println("Please Choose the IP You Want To TraceRoute")
		for i, ip := range ipSlice {
			fmt.Fprintf(color.Output, "%s %s\n",
				color.New(color.FgHiYellow, color.Bold).Sprintf("%d.", i),
				color.New(color.FgWhite, color.Bold).Sprintf("%s", ip),
			)
		}
		var index int
		fmt.Printf("Your Option: ")
		fmt.Scanln(&index)
		if index >= len(ipSlice) || index < 0 {
			fmt.Println("Your Option is invalid")
			os.Exit(3)
		}
		return ipSlice[index]
	}
}
