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
	req := requests.Requests("get").SetDebug(true)
	cookie := &http.Cookie{
		Name:  "anewcookie",
		Value: "20180825",
		Path:  "/",
	}
	resp, err := req.SetCookie(cookie).Run("https://www.httpbin.org")
	if err == nil {
		coo := resp.Cookies()
		// coo is [] *http.Cookies
		println("********cookies*******")
		for _, c := range coo {
			fmt.Println(c.Name, c.Value)
		}
	}

}
