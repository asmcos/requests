package requests

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	// example 1
	println("Get example1")
	req := Requests()

	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Get("http://go.xiulian.net.cn", Header{"Referer": "http://www.jeapedu.com"}, Params{"c": "d", "e": "f"}, Params{"c": "a"})

	// example 2
	println("Get example2")
	h := Header{
		"Referer":         "http://www.jeapedu.com",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
	}

	Get("http://jeapedu.com", h, Header{"accept-encoding": "gzip, deflate, br"})

	// example 3
	println("Get example3")
	p := Params{
		"title": "The blog",
		"name":  "file",
		"id":    "12345",
	}
	resp, err := Requests().Get("http://www.cpython.org", p)

	if err == nil {
		resp.Text()
		fmt.Println(resp.Text())
	}

	// example 4
	println("Get example4")
	// test authentication usernae,password
	//documentation https://www.httpwatch.com/httpgallery/authentication/#showExample10
	req = Requests()
	resp, err = req.Get("https://www.httpwatch.com/httpgallery/authentication/authenticatedimage/default.aspx?0.45874470316137206", Auth{"httpwatch", "foo"})
	if err == nil {
		fmt.Println(resp.R)
	}
	// this save file test PASS
	// resp.SaveFile("auth.jpeg")

	//example 5 test Json
	println("Get example5")
	req = Requests()
	req.Header.Set("Content-Type", "application/json")
	resp, err = req.Get("https://httpbin.org/json")

	if err == nil {
		var json map[string]interface{}
		resp.Json(&json)

		for k, v := range json {
			fmt.Println(k, v)
		}
	}

	// example 6 test gzip
	println("Get example6")
	req = Requests()
	req.Debug = 1
	resp, err = req.Get("https://httpbin.org/gzip")
	if err == nil {
		fmt.Println(resp.Text())
	}
	// example 7 proxy and debug
	println("Get example7")
	req = Requests()
	req.Debug = 1

	// You need open the line
	//req.Proxy("http://192.168.1.190:8888")

	req.Get("https://www.sina.com.cn")

	//example 8 test  auto Cookies
	println("Get example8")
	req = Requests()
	req.Debug = 1
	// req.Proxy("http://192.168.1.190:8888")
	req.Get("https://www.httpbin.org/cookies/set?freeform=1234")
	req.Get("https://www.httpbin.org")
	req.Get("https://www.httpbin.org/cookies/set?a=33d")
	req.Get("https://www.httpbin.org")

	// example 9 test AddCookie
	println("Get example9")
	req = Requests()
	req.Debug = 1

	cookie := &http.Cookie{}
	cookie.Name = "anewcookie"
	cookie.Value = "20180825"
	cookie.Path = "/"

	req.SetCookie(cookie)

	fmt.Println(req.Cookies)
	// req.Proxy("http://127.0.0.1:8888")
	req.Get("https://www.httpbin.org/cookies/set?freeform=1234")
	req.Get("https://www.httpbin.org")
	req.Get("https://www.httpbin.org/cookies/set?a=33d")
	resp, err = req.Get("https://www.httpbin.org")

	if err == nil {
		coo := resp.Cookies()
		// coo is [] *http.Cookies
		println("********cookies*******")
		for _, c := range coo {
			fmt.Println(c.Name, c.Value)
		}
	}
}

func TestPost(t *testing.T) {

	// example 1
	// set post formdata
	println("Post example1")
	req := Requests()
	req.Debug = 1

	data := Datas{
		"comments":  "ew",
		"custemail": "a@231.com",
		"custname":  "1",
		"custtel":   "2",
		"delivery":  "12:45",
		"size":      "small",
		"topping":   "bacon",
	}

	resp, err := req.Post("https://www.httpbin.org/post", data)
	if err == nil {
		fmt.Println(resp.Text())
	}

	//example 2 upload files
	println("Post example2")
	req = Requests()
	req.Debug = 1
	path, _ := os.Getwd()
	path1 := path + "/README.md"
	path2 := path + "/docs/index.md"

	resp, err = req.Post("https://www.httpbin.org/post", data, Files{"a": path1, "b": path2})
	if err == nil {
		fmt.Println(resp.Text())
	}
}

func TestTimeout(t *testing.T) {
	println("Timeout example1")
	req := Requests()
	req.Debug = 1

	// 20 Second
	req.SetTimeout(20)
	req.Get("http://golang.org")

}
