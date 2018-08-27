package main

import "github.com/asmcos/requests"


func main (){

        req := requests.Requests()
        resp := req.Get("https://api.github.com/user",requests.Auth{"asmcos","password...."})
        println(resp.Text())
}
