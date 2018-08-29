package requests

import (
	"fmt"
	"net/http"
	"testing"
	"os"
)



func TestGet(t *testing.T) {
	// example 1
	req := Requests()

	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Get("http://go.xiulian.net.cn", Header{"Content-Length": "0"}, Params{"c": "d", "e": "f"}, Params{"c": "a"})

	// example 2
	h := Header{
		"Referer":         "http://www.jeapedu.com",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
	}

	Get("http://jeapedu.com", h, Header{"Content-Length": "1024"})

	// example 3
	p := Params{
		"title": "The blog",
		"name":  "file",
		"id":    "12345",
	}
	resp := Requests().Get("http://www.cpython.org", p)

	resp.Text()
	fmt.Println(resp.Text())

  // example 4
  // test authentication usernae,password
	//documentation https://www.httpwatch.com/httpgallery/authentication/#showExample10
	req = Requests()
	resp = req.Get("https://www.httpwatch.com/httpgallery/authentication/authenticatedimage/default.aspx?0.45874470316137206",Auth{"httpwatch","foo"})
	fmt.Println(resp.httpresp)

  // this save file test PASS
	// resp.SaveFile("auth.jpeg")

	//example 5 test Json
	req = Requests()
	req.Header.Set("Content-Type","application/json")
	resp = req.Get("https://httpbin.org/json")

	var json map[string]interface{}
  resp.Json(&json)

	for k,v := range json{
		fmt.Println(k,v)
	}

 // example 6 test gzip
 req = Requests()
 req.Debug = 1
 resp = req.Get("https://httpbin.org/gzip")

 fmt.Println(resp.Text())

 // example 7 proxy and debug
 req = Requests()
 req.Debug = 1
 //req.Proxy("http://192.168.1.190:8888")

 resp = req.Get("https://www.sina.com.cn")
 // fmt.Println(resp.Text())
 req.Get("https://www.sina.com.cn")

 //example 8 test  auto Cookies
 req = Requests()
 req.Debug = 1
 // req.Proxy("http://192.168.1.190:8888")
 req.Get("https://www.httpbin.org/cookies/set?freeform=1234")
 req.Get("https://www.httpbin.org")
 req.Get("https://www.httpbin.org/cookies/set?a=33d")
 req.Get("https://www.httpbin.org")

  // example 9 test AddCookie
	req = Requests()
	req.Debug = 1

	cookie := &http.Cookie{}
  cookie.Name   = "anewcookie"
  cookie.Value  = "20180825"
	cookie.Path   = "/"

	req.SetCookie(cookie)


  fmt.Println(req.Cookies)
  // req.Proxy("http://127.0.0.1:8888")
	req.Get("https://www.httpbin.org/cookies/set?freeform=1234")
	req.Get("https://www.httpbin.org")
	req.Get("https://www.httpbin.org/cookies/set?a=33d")
	req.Get("https://www.httpbin.org")

}


func TestPost(t *testing.T) {

  fmt.Println(&http.Cookie{})

  // example 1
	// set post formdata
	req := Requests()
	req.Debug = 1


	data := Datas{
	    "comments": "ew",
	    "custemail": "a@231.com",
	    "custname": "1",
	    "custtel": "2",
	    "delivery": "12:45",
	    "size": "small",
	    "topping": "bacon",
	  }

	resp := req.Post("https://www.httpbin.org/post",data)

	fmt.Println(resp.Text())


  //example 2 upload files
	req = Requests()
	req.Debug = 1
	path, _ := os.Getwd()
  path1 := path +  "/README.md"
	path2 := path +  "/docs/index.md"

  resp = req.Post("https://www.httpbin.org/post",data,Files{"a":path1,"b":path2})

	fmt.Println(resp.Text())

}


func TestTimeout(t *testing.T) {
	req := Requests()
	req.Debug = 1

  // 20 Second
	req.SetTimeout(20)
	req.Get("http://golang.org")

}
