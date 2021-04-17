package examples

import (
	"fmt"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

func TestPostJson(t *testing.T) {
	println("Test POST: post data and json")
	data := requests.Datas{
		"comments": "ew",
	}
	json := requests.Json{
		"key": "value",
	}
	json = map[string]interface{}{
		"key": "value",
	}
	resp, err := requests.Post("https://www.httpbin.org/post", data, json)
	if err == nil {
		fmt.Println(resp.Text())
	}
}
func TestPostString(t *testing.T) {
	println("Test POST: post data and json")
	dataStr := "{\"key\":\"This is raw data\"}"
	resp, err := requests.Post("https://www.httpbin.org/post", dataStr)
	if err == nil {
		fmt.Println(resp.Text())
	}
}
