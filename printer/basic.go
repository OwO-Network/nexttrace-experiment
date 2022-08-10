package printer

import (
	"fmt"
	"net"

	"github.com/fatih/color"
)

var version = "v0.0.0.alpha"
var buildDate = ""
var commitID = ""

func Version() {
	fmt.Fprintf(color.Output, "%s %s %s %s\n",
		color.New(color.FgWhite, color.Bold).Sprintf("%s", "NextTrace Enhanced"),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", version),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", buildDate),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", commitID),
	)
}

func CopyRight() {
	fmt.Fprintf(color.Output, "%s\n%s %s\n%s %s\n%s %s\n%s %s\n\n",
		color.New(color.FgGreen, color.Bold).Sprintf("%s", "NextTrace Project Contributor"),
		color.New(color.FgWhite, color.Bold).Sprintf("%s", "Tso"),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", "tsosunchia@gmail.com"),
		color.New(color.FgWhite, color.Bold).Sprintf("%s", "Vincent"),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", "vincent.moe"),
		color.New(color.FgWhite, color.Bold).Sprintf("%s", "zhshch"),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", "xzhsh.ch"),
		color.New(color.FgWhite, color.Bold).Sprintf("%s", "Leo"),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", "leo.moe"),
	)

	PluginCopyRight()
}

func PluginCopyRight() {
	fmt.Fprintf(color.Output, "%s\n%s %s\n",
		color.New(color.FgGreen, color.Bold).Sprintf("%s", "NextTrace Enhanced Map Plugin"),
		color.New(color.FgWhite, color.Bold).Sprintf("%s", "Tso"),
		color.New(color.FgHiBlack, color.Bold).Sprintf("%s", "tsosunchia@gmail.com"),
	)
}

func PrintTraceRouteNav(ip net.IP, domain string, dataOrigin string) {
	fmt.Println("IP Geo Data Provider: " + dataOrigin)

	if ip.String() == domain {
		fmt.Printf("traceroute to %s, 30 hops max, 32 byte packets\n", ip.String())
	} else {
		fmt.Printf("traceroute to %s (%s), 30 hops max, 32 byte packets\n", ip.String(), domain)
	}
}
