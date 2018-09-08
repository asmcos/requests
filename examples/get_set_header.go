package main

import "github.com/asmcos/requests"


func main (){

        req := requests.Requests()

        resp,_ := req.Get("http://go.xiulian.net.cn",requests.Header{"Referer":"http://www.jeapedu.com"})

        println(resp.Text())

}
