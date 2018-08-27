# requests

Requests is an HTTP library  , it is easy to use. Similar to Python requests.

# Installation

```
go get -u github.com/asmcos/requests
```

# Start

``` go
package main

import "github.com/asmcos/requests"

func main (){

        resp := requests.Get("http://go.xiulian.net.cn")
        println(resp.Text())
}
```

## Post

``` go
package main

import "github.com/asmcos/requests"


func main (){

        data := requests.Datas{
          "name":"requests_post_test",
        }
        resp := requests.Post("https://www.httpbin.org/post",data)
        println(resp.Text())
}

```

     Server return data...

``` json
{
  "args": {},
  "data": "",
  "files": {},
  "form": {
    "name": "requests_post_test"
  },
  "headers": {
    "Accept-Encoding": "gzip",
    "Connection": "close",
    "Content-Length": "23",
    "Content-Type": "application/x-www-form-urlencoded",
    "Host": "www.httpbin.org",
    "User-Agent": "Go-Requests 0.5"
  },
  "json": null,
  "origin": "114.242.34.110",
  "url": "https://www.httpbin.org/post"
}

```

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


# Set header

### exmaple 1

``` go
req := requests.Requests()

resp := req.Get("http://go.xiulian.net.cn",requests.Header{"Referer":"http://www.jeapedu.com"})
println(resp.Text())
```

### example 2

``` go
req := requests.Requests()
req.Header.Set("accept-encoding", "gzip, deflate, br")
resp := req.Get("http://go.xiulian.net.cn",requests.Header{"Referer":"http://www.jeapedu.com"})
println(resp.Text())

```

### example 3

``` go
h := requests.Header{
  "Referer":         "http://www.jeapedu.com",
  "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
}
resp := req.Get("http://go.xiulian.net.cn",h)

h2 := requests.Header{
  ...
  ...
}
h3,h4 ....
// two or more headers ...
resp = req.Get("http://go.xiulian.net.cn",h,h2,h3,h4)
```


# Set params

``` go
p := requests.Params{
  "title": "The blog",
  "name":  "file",
  "id":    "12345",
}
resp := req.Get("http://www.cpython.org", p)

```


# Auth

Test with the `correct` user information.

``` go
req := requests.Requests()
resp := req.Get("https://api.github.com/user",requests.Auth{"asmcos","password...."})
println(resp.Text())
```

github return

```
{"login":"asmcos","id":xxxxx,"node_id":"Mxxxxxxxxx==".....
```

# JSON

``` go
req := requests.Requests()
req.Header.Set("Content-Type","application/json")
resp = req.Get("https://httpbin.org/json")

var json map[string]interface{}
resp.Json(&json)

for k,v := range json{
  fmt.Println(k,v)
}
```
