package examples

import (
	"testing"

	"github.com/asmcos/requests"
	_ "github.com/asmcos/requests/init"
)

func TestProxy(t *testing.T) {
	println("5. Get: with proxy")
	session := requests.Sessions()
	// session.Proxy("http://192.168.1.190:8888")
	session.Get("https://www.httpbin.org/cookies/set?freeform=1234")
	session.Get("https://www.httpbin.org")
}
