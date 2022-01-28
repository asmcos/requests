package examples

import (
	ejson "encoding/json"
	"testing"

	"github.com/asmcos/requests"
	_ "github.com/asmcos/requests/init"
)

// Post params
func TestPostParams(t *testing.T) {
	println("Test POST: post params")
	resp, err := requests.Post(
		"https://www.httpbin.org/post",
		requests.Params{
			"name": "asmcos",
		},
	)
	if err != nil {
		t.Error(err)
	}
	var data = struct {
		Args struct {
			Name string
		}
	}{}
	_ = resp.Json(&data)
	if data.Args.Name != "asmcos" {
		t.Error("invalid response body:", resp.Text())
	}
}

// Post Form Data
func TestPostForm(t *testing.T) {
	println("Test POST: post form data(x-wwww-form-urlencoded)")
	resp, err := requests.Post(
		"https://www.httpbin.org/post",
		requests.Datas{
			"name": "asmcos",
		},
	)
	if err != nil {
		t.Error(err)
	}
	var data = struct {
		Form struct {
			Name string
		}
	}{}
	err = resp.Json(&data)
	if data.Form.Name != "asmcos" {
		t.Error("invalid response body:", resp.Text())
	}
}

// Post Json Request
func TestPostJson(t *testing.T) {
	println("Test POST: post json data")
	json := requests.Json{
		"name": "Alex",
	}
	/*
		    //it still works!
			json = map[string]interface{}{
				"key": "value",
			}
	*/
	resp, err := requests.Post("https://www.httpbin.org/post", json)
	if err != nil {
		t.Error(err)
	}

	// parse data
	var data = struct {
		Data string
	}{}
	resp.Json(&data)

	// is expected results
	jsonData, _ := ejson.Marshal(json) // if data.Data!= "{\"name\":\"Alex\"}"{
	if data.Data != string(jsonData) {
		t.Error("invalid response body:", resp.Text())
	}
}

// Post QueryString: application/x-www-form-urlencoded
func TestPostQueryString(t *testing.T) {
	println("Test POST: raw post data ")
	queryString := "name=Alex&age=29"
	resp, err := requests.Post("https://www.httpbin.org/post", queryString)
	if err != nil {
		t.Fatal(err)
	}
	var data = struct {
		Form struct {
			Name string
			Age  string
		}
	}{}
	err = resp.Json(&data)
	if data.Form.Age != "29" {
		t.Error("invalid response body:", resp.Text())
	}
}

// Post Raw text/plain
func TestRawString(t *testing.T) {
	println("Test POST: raw post data ")
	rawText := "raw data: Hi, Jack!"
	resp, err := requests.Post("https://www.httpbin.org/post", rawText,
		requests.Header{"Content-Type": "text/plain"},
	)
	if err != nil {
		t.Fatal(err)
	}
	var data interface{}
	err = resp.Json(&data)
	if data.(map[string]interface{})["data"].(string) != rawText {
		t.Error("invalid response body:", resp.Text())
	}
}

// Post Raw text/plain (with bytes)
func TestRawBytes(t *testing.T) {
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
		t.Error("invalid response body:", resp.Text())
	}

}
