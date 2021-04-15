package main

import "github.com/ahuigo/requests"


func main (){

        resp,_ := requests.Get("https://www.baidu.com/")
        println(resp.Text())
        println(resp.R.Header.Get("location"))
        println(resp.R.Header.Get("Location"))
}
