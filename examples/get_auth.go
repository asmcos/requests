package main

import (
	"github.com/ahuigo/requests"
	"fmt"
)

func main (){

	req := requests.Requests()

  resp,_ := req.Get("https://api.github.com/user",requests.Auth{"ahuigo","password...."})
  println(resp.Text())
	fmt.Println(resp.R.StatusCode)
	fmt.Println(resp.R.Header["Content-Type"])
}
