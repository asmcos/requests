package main

import (
	"github.com/asmcos/requests"
	"fmt"
)

func main (){

	req := requests.Requests()

  resp,_ := req.Get("https://api.github.com/user",requests.Auth{"asmcos","password...."})
  println(resp.Text())
	fmt.Println(resp.R.StatusCode)
	fmt.Println(resp.R.Header["Content-Type"])
}
