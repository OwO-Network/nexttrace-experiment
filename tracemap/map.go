package tracemap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetMapUrl(r string) {
	url := "https://api.leo.moe/tracemap/api"
	resp, _ := http.Post(url, "application/json", strings.NewReader(string(r)))
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("MapTrace URL: " + string(body))
}
