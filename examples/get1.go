package main

import "github.com/asmcos/requests"


func main (){

        resp,_ := requests.Get("http://go.xiulian.net.cn")
        println(resp.Text())
}
