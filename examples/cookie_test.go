package examples

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

func TestSendCookie(t *testing.T) {
	// example 9 test AddCookie
	println("6. Get: get with cookie")
	req := requests.Requests().SetDebug(true)
	cookie := &http.Cookie{
		Name:  "anewcookie",
		Value: "20180825",
		Path:  "/",
	}
	resp, err := req.SetCookie(cookie).Get("https://www.httpbin.org")
	// resp, err := req.SetCookie(cookie).Get("http://127.0.0.1:8088")
	if err != nil {
		panic(err)
	}
	coo := resp.Cookies()
	// coo is [] *http.Cookies
	println("********cookies*******")
	for _, c := range coo {
		fmt.Println(c.Name, c.Value)
	}

}
