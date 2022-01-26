package examples

import (
	"testing"

	"github.com/ahuigo/requests"
)

// Send headers
func TestSendHeaders(t *testing.T) {
	println("Test Get: send header")
	requests.Get(
		"http://www.zhanluejia.net.cn",
		requests.Header{"Referer": "http://www.jeapedu.com"},
	)
}

// Set session headers
func TestSendSessionHeader(t *testing.T) {
	session := requests.Sessions()
	session.SetHeader("accept-encoding", "gzip, deflate, br")
	session.Get("http://www.zhanluejia.net.cn")
}

// Set global header(user-agent)
func TestSetGlobalHeader(t *testing.T) {
	headerK := "User-Agent"
	headerV := "Custom-Test-Go-User-Agent"
	requests.SetHeader(headerK, headerV)
	req, err := requests.BuildRequest("post", "http://baidu.com/a/b/c")
	if err != nil {
		t.Fatal(err)
	}
	if req.Header.Get(headerK) != headerV {
		t.Fatalf("Expected header %s is %s", headerK, headerV)
	}
}
