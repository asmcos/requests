# Requests
[![license](http://dmlc.github.io/img/apache2.svg)](https://raw.githubusercontent.com/ahuigo/requests/master/LICENSE)

# Requests
Requests is an HTTP library, it is easy to use. Similar to Python requests.

Warning: Session is not safe in multi goroutine. You can not do as following:

    // Bad! Do not call session in in multi goroutine!!!!!
    session := requests.Sessions()

    // goroutine 1
    go func(){
       session.Post(url1) 
    }()

    // goroutine 2
    go func(){
       session.Post(url2) 
    }()


# Installation

```
go get -u github.com/ahuigo/requests
```

# Examples
> For more examples, refer to https://github.com/ahuigo/requests/tree/master/examples

## Get
    package main
    import (
        "github.com/ahuigo/requests"
        "fmt"
    )

    func main(){
        var json map[string]interface{}
        params := requests.Params{"name": "ahuigo", "page":"1"}
        resp, err := requests.Get("https://httpbin.org/json", params)
        if err == nil {
            resp.Json(&json)
            for k, v := range json {
                fmt.Println(k, v)
            }
        }
    }


## Post

### Post params

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

### Post Form Data
    // Post Form Data
    func TestPostForm(t *testing.T) {
        println("Test POST: post form data")
        data := requests.Datas{
            "comments": "ew",
        }
        resp, err := requests.Post("https://www.httpbin.org/post", data)
        if err == nil {
            fmt.Println(resp.Text())
        }
    }


### Post Json 
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

### Post Raw Text
    func TestPostString(t *testing.T) {
        println("Test POST: post data and json")
        rawText := "raw data: Hi, Jack!"
        resp, err := requests.Post("https://www.httpbin.org/post", rawText,
            requests.Header{"Content-Type": "text/plain"},
        )
        if err == nil {
            fmt.Println(resp.Text())
        }
    }

### PostFiles

	path, _ := os.Getwd()
	session := requests.Sessions()

	resp, err := session.SetDebug(true).Post(
		"https://www.httpbin.org/post",
		requests.Files{
            "file1": path + "/README.md",
            "file2": path + "/version",
        },
	)
	if err == nil {
		fmt.Println(resp.Text())
	}

## Session Support
    // 0. Make a session
	session := r.Sessions()

    // 1. First, set cookies: count=100
	var data struct {
		Cookies struct {
			Count string `json:"count"`
		}
	}
	session.Get("https://httpbin.org/cookies/set?count=100")

	// 2. Second, get cookies
	resp, err := session.Get("https://httpbin.org/cookies")
	if err == nil {
		resp.Json(&data)
        if data.Cookies.Count!="100"{
            t.Fatal("Failed to get valid cookies: "+resp.Text())
        }
	}

Warning: Session is not safe in multi goroutine. You can not do as following:

    // Bad! Do not call session in in multi goroutine!!!!!
    session := requests.Sessions()

    // goroutine 1
    go func(){
       session.Post(url1) 
    }()

    // goroutine 2
    go func(){
       session.Post(url2) 
    }()

## Request Options

### SetTimeout

    session := Requests.Sessions()
    session.SetTimeout(20)

### Debug Mode
    
    session.Debug = 1

### Set Authentication
    session := requests.Sessions()
    resp,_ := session.Get("https://api.github.com/user",requests.Auth{"asmcos","password...."})

### Set Cookie
	cookie1 := http.Cookie{Name: "cookie_name", Value: "cookie_value"}
    session.SetCookie(&cookie1)

### Set header

    func TestGetParamsHeaders(t *testing.T) {
        println("Test Get: custom header and params")
        requests.Get("http://www.zhanluejia.net.cn",
            requests.Header{"Referer": "http://www.jeapedu.com"},
            requests.Params{"page": "1", "size": "20"},
            requests.Params{"name": "ahuio"},
        )
    }

    func TestGetParamsHeaders2(t *testing.T) {
        session := requests.Sessions()
        session.SetHeader("accept-encoding", "gzip, deflate, br")
        session.Run("http://www.zhanluejia.net.cn",
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

two or more headers ...

    headers1 := requests.Header{"Referer": "http://www.jeapedu.com"},
    ....
    resp,_ = session.Get(
        "http://www.zhanluejia.net.cn",
        headers1,
        headers2,
        headers3,
    )


## Response
### Fetch Response Body
https://github.com/ahuigo/requests/blob/master/examples/resp_test.go

    fmt.Println(resp.Text())
    fmt.Println(resp.Content())

### Fetch Response Cookies
https://github.com/ahuigo/requests/blob/master/examples/cookie_test.go

    resp,_ = session.Get("https://www.httpbin.org")
    coo := resp.Cookies()
    for _, c:= range coo{
        fmt.Println(c.Name,c.Value)
    }

## Build Request

	req, err := r.BuildRequest("post", "http://baidu.com/a/b/c", r.Json{
		"age": 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(req.Body)
	expectedBody := `{"age":1}`
	if string(body) != expectedBody {
		t.Fatal("Failed to build request")
	}

## Custom

### Custom User Agent

	headerK := "User-Agent"
	headerV := "Custom-Test-Go-User-Agent"
	requests.SetHeader(headerK, headerV)

# Feature Support
  - Set headers
  - Set params
  - Multipart File Uploads
  - Sessions with Cookie Persistence
  - Proxy
  - Authentication
  - JSON
  - Chunked Requests
  - Debug
  - SetTimeout

# Thanks
This project is inspired by [github.com/asmcos/requests](http://github.com/asmcos/requests). 

Great thanks to it :).
