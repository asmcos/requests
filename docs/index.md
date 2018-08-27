# Installation


# example 1

```
package main

import "github.com/asmcos/requests"


func main (){

        resp := requests.Get("http://go.xiulian.net.cn")
        println(resp.Text())
}

```
