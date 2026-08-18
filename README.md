
[![license](http://dmlc.github.io/img/apache2.svg)](https://raw.githubusercontent.com/asmcos/requests/master/LICENSE)

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

        resp,err := requests.Get("http://www.zhanluejia.net.cn")
        if err != nil{
          return
        }
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
        resp,_ := requests.Post("https://www.httpbin.org/post",data)
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


## PostJson

``` go
package main

import "github.com/asmcos/requests"


func main (){

        jsonStr := "{\"name\":\"requests_post_test\"}"
        resp,_ := requests.PostJson("https://www.httpbin.org/post",jsonStr)
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
  - SetTimeout
  - Skip TLS certificate verification (`SetInsecureSkipVerify`)

# Changelog (v0.9)

This release was completed with AI coding assistance. The following bugs from GitHub issues and pull requests were fixed:

- **[#33](https://github.com/asmcos/requests/issues/33) / [PR #38](https://github.com/asmcos/requests/pull/38)** — HTTP/2 retry failed with `cannot retry err ... define Request.GetBody`. Request bodies are now rewindable via `GetBody`.
- **[#37](https://github.com/asmcos/requests/issues/37)** — `PostJson` sent extra bytes in the body (chunked encoding / trailing newline from `json.Encoder`). Bodies now use `Content-Length` and `json.Marshal`.
- **[#10](https://github.com/asmcos/requests/issues/10)** — `http: ContentLength=N with Body length 0` when calling `Get` after `Post` on the same client. Stale body and `Content-Length` are cleared between requests.
- **[PR #36](https://github.com/asmcos/requests/pull/36)** — Multipart uploads wrote files before form fields, which broke S3-style POSTs that require `key` first. Form fields are now written before files.
- **[PR #39](https://github.com/asmcos/requests/pull/39)** — Errors were printed with `fmt.Println`. Errors are returned only, without extra stdout.
- **[#23](https://github.com/asmcos/requests/issues/23) / [#7](https://github.com/asmcos/requests/issues/7)** — `Post` always overwrote `Content-Type`. The default `application/x-www-form-urlencoded` is set only when the caller did not set one.
- **[#28](https://github.com/asmcos/requests/issues/28)** — No way to skip TLS verification (`verify=False`). Added `SetInsecureSkipVerify`. `Proxy()` no longer forces `InsecureSkipVerify`.
- **[PR #22](https://github.com/asmcos/requests/pull/22)** — `PostJson` ignored URL query `Params`. Query parameters are now applied.
- **Host header** — `Header{"Host": "..."}` did not change the HTTP Host. It is now applied to `Request.Host`.
- **Resource leaks** — Uploaded files and gzip readers were not closed. Both are closed after use.


# Set header

### example 1

``` go
req := requests.Requests()

resp,err := req.Get("http://www.zhanluejia.net.cn",requests.Header{"Referer":"http://www.jeapedu.com"})
if (err == nil){
  println(resp.Text())
}
```

### example 2

``` go
req := requests.Requests()
req.Header.Set("accept-encoding", "gzip, deflate, br")
resp,_ := req.Get("http://www.zhanluejia.net.cn",requests.Header{"Referer":"http://www.jeapedu.com"})
println(resp.Text())

```

### example 3

``` go
h := requests.Header{
  "Referer":         "http://www.jeapedu.com",
  "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
}
resp,_ := req.Get("http://wwww.zhanluejia.net.cn",h)

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
resp,_ := req.Get("http://www.cpython.org", p)

```


# Auth

Test with the `correct` user information.

``` go
req := requests.Requests()
resp,_ := req.Get("https://api.github.com/user",requests.Auth{"asmcos","password...."})
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
resp,_ = req.Get("https://httpbin.org/json")

var json map[string]interface{}
resp.Json(&json)

for k,v := range json{
  fmt.Println(k,v)
}
```


# SetTimeout

```
req := Requests()
req.Debug = 1

// 20 Second
req.SetTimeout(20)
req.Get("http://golang.org")
```

# Get Cookies

``` go
resp,_ = req.Get("https://www.httpbin.org")
coo := resp.Cookies()
// coo is [] *http.Cookies
println("********cookies*******")
for _, c:= range coo{
  fmt.Println(c.Name,c.Value)
}
```
