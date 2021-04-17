
[![license](http://dmlc.github.io/img/apache2.svg)](https://raw.githubusercontent.com/ahuigo/requests/master/LICENSE)

# requests

Requests is an HTTP library  , it is easy to use. Similar to Python requests.

# Installation

```
go get -u github.com/ahuigo/requests
```

# Start

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
    // json := requests.Json{ "key": "value"}
    json = map[string]interface{}{
        "key": "value",
    }
    resp, err := requests.Post("https://www.httpbin.org/post", data, json)
    if err == nil {
        fmt.Println(resp.Text())
    }

### PostString

    dataStr := "{\"key\":\"This is raw data\"}"
    resp, err := requests.Post("https://www.httpbin.org/post", dataStr)
    if err == nil {
        fmt.Println(resp.Text())
    }

### PostFiles

	path, _ := os.Getwd()
	req := requests.Requests("GET").SetDebug(true)

	resp, err := req.SetMethod("POST").Run(
		"https://www.httpbin.org/post",
		requests.Files{
            "file1": path + "/README.md",
            "file2": path + "/version",
        },
	)
	if err == nil {
		fmt.Println(resp.Text())
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


# Set header

### example 1

``` go
resp,err := requests.Get("http://www.zhanluejia.net.cn",requests.Header{"Referer":"http://www.jeapedu.com"})
if (err == nil){
  println(resp.Text())
}
```

### example 2

``` go
req := requests.Requests("GET")
req.Header.Set("accept-encoding", "gzip, deflate, br")
resp,_ := req.Run("http://www.zhanluejia.net.cn",requests.Header{"Referer":"http://www.jeapedu.com"})
println(resp.Text())

```

### example 3

``` go
h := requests.Header{
  "Referer":         "http://www.jeapedu.com",
  "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
}
resp,_ := requests.Get("http://wwww.zhanluejia.net.cn",h)

h2 := requests.Header{
  ...
  ...
}
h3,h4 ....
// two or more headers ...
resp,_ = req.Get("http://www.zhanluejia.net.cn",h,h2,h3,h4)
```


# Set params

``` go
p := requests.Params{
  "title": "The blog",
  "name":  "file",
  "id":    "12345",
}
resp,_ := requests.Get("http://www.cpython.org", p)

```


# Auth

Test with the `correct` user information.

``` go
resp,_ := requests.Get("https://api.github.com/user",requests.Auth{"ahuigo","password...."})
println(resp.Text())
```

github return

```
{"login":"ahuigo","id":xxxxx,"node_id":"Mxxxxxxxxx==".....
```

# JSON

``` go
req.Header.Set("Content-Type","application/json")
resp,_ = requests.Get("https://httpbin.org/json")

var json map[string]interface{}
resp.Json(&json)

for k,v := range json{
  fmt.Println(k,v)
}
```


# SetTimeout

```
req := Requests("GET")
req.Debug = 1

// 20 Second
req.SetTimeout(20)
req.Run("http://golang.org")
```

# Get Cookies

``` go
resp,_ = requests.Get("https://www.httpbin.org")
coo := resp.Cookies()
// coo is [] *http.Cookies
println("********cookies*******")
for _, c:= range coo{
  fmt.Println(c.Name,c.Value)
}
```

# Thanks
This project is inspired by [github.com/asmcos/requests](http://github.com/asmcos/requests). 

Great thanks to it :).
