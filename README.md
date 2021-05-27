# Requests
[![license](http://dmlc.github.io/img/apache2.svg)](https://raw.githubusercontent.com/ahuigo/requests/master/LICENSE)

# requests
Requests is an HTTP library, it is easy to use. Similar to Python requests.

# Installation

```
go get -u github.com/ahuigo/requests
```

# Examples
> For more examples, refer to https://github.com/ahuigo/requests/tree/master/examples

## Get

    var json map[string]interface{}
    resp, err := requests.Get("https://httpbin.org/json")
    if err == nil {
        resp.Json(&json)
        for k, v := range json {
            fmt.Println(k, v)
        }
    }

## Post

### PostJson
    data := requests.Datas{
        "comments": "ew",
    }
    json := map[string]interface{}{
        "key": "value",
    }
    resp, err := requests.Post("https://www.httpbin.org/post", data, json)
    if err == nil {
        fmt.Println(resp.Text())
    }

You can use json builder instead:

    json := requests.Json{ "key": "value"}

### PostFiles

	path, _ := os.Getwd()
	req := requests.Requests()

	resp, err := req.SetDebug(true).Post(
		"https://www.httpbin.org/post",
		requests.Files{
            "file1": path + "/README.md",
            "file2": path + "/version",
        },
	)
	if err == nil {
		fmt.Println(resp.Text())
	}

## Request Options

### SetTimeout

    req := Requests()
    req.SetTimeout(20)


### Debug Mode
    
    req.Debug = 1

### Set Authentication
    req := requests.Requests()
    resp,_ := req.Get("https://api.github.com/user",requests.Auth{"asmcos","password...."})

### Set Cookie
	cookie1 := http.Cookie{Name: "cookie_name", Value: "cookie_value"}
    req.SetCookie(&cookie1)

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

two or more headers ...

    headers1 := requests.Header{"Referer": "http://www.jeapedu.com"},
    ....
    resp,_ = req.Get(
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

    resp,_ = req.Get("https://www.httpbin.org")
    coo := resp.Cookies()
    for _, c:= range coo{
        fmt.Println(c.Name,c.Value)
    }


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
