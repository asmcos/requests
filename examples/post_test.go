package examples

import (
	"fmt"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

// Post params
func TestPostParams(t *testing.T) {
	println("Test POST: post params")
	data := requests.Params{
		"name": "ahuigo",
	}
	resp, err := requests.Post("https://www.httpbin.org/post", data)
	if err == nil {
		fmt.Println(resp.Text())
	}
}

// Post Form Request
func TestPostForm(t *testing.T) {
	println("Test POST: post form data(x-wwww-form-urlencoded)")
	data := requests.Datas{
		"comments": "ew",
	}
	resp, err := requests.Post("https://www.httpbin.org/post", data)
	if err == nil {
		fmt.Println(resp.Text())
	}
}

// Post Json Request
func TestPostJson(t *testing.T) {
	println("Test POST: post json data")
	json := requests.Json{
		"key": "value",
	}
	/*
		    //it still works!
			json = map[string]interface{}{
				"key": "value",
			}
	*/
	resp, err := requests.Post("https://www.httpbin.org/post", json)
	if err == nil {
		fmt.Println(resp.Text())
	}
}

// Post Raw Text
func TestPostString(t *testing.T) {
	println("Test POST: raw post data ")
	rawText := "raw data: Hi, Jack!"
	resp, err := requests.Post("https://www.httpbin.org/post", rawText,
		requests.Header{"Content-Type": "text/plain"},
	)
	if err == nil {
		fmt.Println(resp.Text())
	}
}

// Post Raw Text
func TestPostBytes(t *testing.T) {
	println("Test POST: post bytes data")
	rawText := "raw data: Hi, Jack!"
	resp, err := requests.Post("https://www.httpbin.org/post", []byte(rawText),
		requests.Header{"Content-Type": "text/plain"},
	)
	if err != nil {
		t.Error(err)
	}
	var data = struct {
		Data string
	}{}
	err = resp.Json(&data)
	if data.Data != rawText {
		t.Error(err)
	}

}
