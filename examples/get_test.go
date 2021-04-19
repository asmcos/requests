package examples

import (
	"fmt"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

func TestGetJson(t *testing.T) {
	println("Test Get: fetch json response")
	resp, err := requests.Get("https://httpbin.org/json")
	if err == nil {
		var json map[string]interface{}
		resp.Json(&json)
		for k, v := range json {
			fmt.Println(k, v)
		}
	}
}

func TestGetParamsHeaders(t *testing.T) {
	println("Test Get: custom header and params")
	requests.Get("http://www.zhanluejia.net.cn",
		requests.Header{"Referer": "http://www.jeapedu.com"},
		requests.Params{"page": "1", "size": "20"},
		requests.Params{"name": "ahuio"},
	)
}
func TestGetParamsHeaders2(t *testing.T) {
	req := requests.Requests("get")
	req.SetHeader("accept-encoding", "gzip, deflate, br")
	req.Run("http://www.zhanluejia.net.cn",
		requests.Params{"page": "1", "size": "20"},
		requests.Params{"name": "ahuio"},
	)
}
func TestResponseHeader(t *testing.T) {
	resp, _ := requests.Get("https://www.baidu.com/")
	println(resp.Text())
	println(resp.R.Header.Get("location"))
	println(resp.R.Header.Get("Location"))
}
