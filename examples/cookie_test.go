package examples

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/asmcos/requests"
	_ "github.com/asmcos/requests/init"
)

func TestSendCookie(t *testing.T) {
	resp, err := requests.Get("https://www.httpbin.org/cookies",
		requests.Header{"Cookie": "id_token=1234"},
		requests.Json{"workflow_id": "wfid1234"},
	)
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{}
	resp.Json(&data)
	fmt.Println(data)

}

// Test session Cookie
func TestSessionCookie(t *testing.T) {
	println("Test: send and get cookie")
	req := requests.Sessions().SetDebug(true)
	cookie := &http.Cookie{
		Name:  "name1",
		Value: "value1",
		Path:  "/",
	}
	req.SetCookie(cookie).Get("https://www.httpbin.org")
	resp, err := req.Get("https://www.httpbin.org",
		&http.Cookie{
			Name:  "name2",
			Value: "value2",
		},
	)
	// resp, err := req.SetCookie(cookie).Get("http://127.0.0.1:8088")
	if err != nil {
		panic(err)
	}
	cookies := map[string]string{}
	// cookies's type is `[]*http.Cookies`
	println("********session cookies*******")
	for _, c := range resp.Cookies() {
		if _, exists := cookies[c.Name]; exists {
			t.Fatal("duplicated cookie:", c.Name, c.Value)
		}
		cookies[c.Name] = c.Value
	}
	if cookies["name1"] != "value1" || cookies["name2"] != "value2" {
		t.Fatal("Failed to send valid cookie")
	}

}
