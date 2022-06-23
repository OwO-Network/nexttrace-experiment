package printer

import (
	"encoding/json"

	"github.com/xgadget-lab/nexttrace/trace"
)

func ParseJson(res *trace.Result) string {
	r, _ := json.Marshal(res)
	return string(r)
}
