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

## example Post

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

     Server return some datas...

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
