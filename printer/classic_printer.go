package printer

import (
	"fmt"

	"github.com/OwO-Network/nexttrace-enhanced/trace"
)

func ClassicPrinter(res *trace.Result, ttl int) {
	fmt.Print(ttl + 1)
	for i := range res.Hops[ttl] {
		HopPrinter(res.Hops[ttl][i])
	}
}
