package printer

import (
	"encoding/json"

	"github.com/OwO-Network/nexttrace-enhanced/trace"
)

func ParseJson(res *trace.Result) string {
	r, _ := json.Marshal(res)
	return string(r)
}
