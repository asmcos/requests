# requests

Requests is an HTTP library  , it is easy to use. Similar to Python requests.

# Install

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

# Set header

### exmaple 1

```
req := requests.Requests()

resp := req.Get("http://go.xiulian.net.cn",requests.Header{"Referer":"http://www.jeapedu.com"})
println(resp.Text())
```

### example 2

```
req := requests.Requests()
req.Header.Set("accept-encoding", "gzip, deflate, br")
resp := req.Get("http://go.xiulian.net.cn",requests.Header{"Referer":"http://www.jeapedu.com"})
println(resp.Text())

```

### example 3

```
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

```
p := requests.Params{
  "title": "The blog",
  "name":  "file",
  "id":    "12345",
}
resp := req.Get("http://www.cpython.org", p)

```
