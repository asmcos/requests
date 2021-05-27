package examples

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

// Test send/get Cookie
func TestCookie(t *testing.T) {
	println("Test: send and get cookie")
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
	cookies := resp.Cookies()
	// cookies's type is `[]*http.Cookies`
	println("********session cookies*******")
	for _, c := range cookies {
		fmt.Println(c.Name, c.Value)
	}

}
