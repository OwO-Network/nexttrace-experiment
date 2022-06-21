package printer

import (
	"encoding/json"

	"github.com/xgadget-lab/nexttrace/trace"
)

func PrintJson(res *trace.Result) {
	r, _ := json.Marshal(res)
	println(string(r))
}
