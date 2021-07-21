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

type HbResponse struct {
	Args map[string]string `json:"args"`
}

// Get with params
func TestGetParams(t *testing.T) {
	params := requests.Params{"name": "ahuigo"}
	resp, err := requests.Get("https://httpbin.org/get", params)
	if err == nil {
		json := &HbResponse{}
		if err := resp.Json(&json); err != nil {
			t.Fatal(err)
		}
		if json.Args["name"] != "ahuigo" {
			t.Fatal("Invalid response: " + resp.Text())
		}
	}
	if err != nil {
		t.Fatal(err)
	}
}

// Send headers
func TestGetHeaders(t *testing.T) {
	println("Test Get: custom header")
	requests.Get(
		"http://www.zhanluejia.net.cn",
		requests.Header{"Referer": "http://www.jeapedu.com"},
	)
}

// Send headers
func TestGetHeaderParams(t *testing.T) {
	session := requests.Sessions()
	session.SetHeader("accept-encoding", "gzip, deflate, br")
	session.Get(
		"http://www.zhanluejia.net.cn",
		requests.Params{"name": "ahuio"},
	)
}
