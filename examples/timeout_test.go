package examples

import (
	"fmt"
	"testing"

	"github.com/asmcos/requests"
	_ "github.com/asmcos/requests/init"
	"github.com/davecgh/go-spew/spew"
)

func TestClose(t *testing.T) {
	fmt.Println("Test Close")
	req := requests.Sessions()
	for i := 0; i < 10; i++ {
		_, err := req.Post(
			"http://localhost:1337/requests",
			requests.Datas{"SrcIp": "4312"})
		fmt.Printf("\r%d %v", i, err)
		req.Close()
	}

	spew.Dump(req)
	fmt.Println("10 times get test end.")
}
func TestTimeout(t *testing.T) {
	println("Test Timeout")
	req := requests.Sessions().SetTimeout(20)
	req.Get("http://golang.org")
}
