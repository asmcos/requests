package examples

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

func TestGetDebug(t *testing.T) {
	println("4. Get: SetDebug")
	session := requests.Sessions().SetDebug(true)
	resp, err := session.Post("https://httpbin.org/post",
		requests.Json{
			"name": "ahuigo",
		},
		&http.Cookie{
			Name:  "count",
			Value: "1",
		},
	)
	if err == nil {
		fmt.Println("response text:", resp.Text())
	}
}
