package examples

import (
	"fmt"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

// Get Json Response
func TestGetJson(t *testing.T) {
	println("Test Get: fetch json response")
	resp, err := requests.Get("https://httpbin.org/json")
	if err == nil {
		var json map[string]interface{}
		err = resp.Json(&json)
		for k, v := range json {
			fmt.Println(k, v)
		}
	}
	if err != nil {
		t.Fatal(err)
	}
}

// Send headers
func TestGetParamsHeaders(t *testing.T) {
	println("Test Get: custom header and params")
	requests.Get("http://www.zhanluejia.net.cn",
		requests.Header{"Referer": "http://www.jeapedu.com"},
		requests.Params{"page": "1", "size": "20"},
		requests.Params{"name": "ahuio"},
	)
}

// Send headers
func TestGetParamsHeaders2(t *testing.T) {
	req := requests.Requests()
	req.SetHeader("accept-encoding", "gzip, deflate, br")
	req.Get("http://www.zhanluejia.net.cn",
		requests.Params{"page": "1", "size": "20"},
		requests.Params{"name": "ahuio"},
	)
}

