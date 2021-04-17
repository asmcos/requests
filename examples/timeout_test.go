package examples

import (
	"fmt"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
	"github.com/davecgh/go-spew/spew"
)

func TestClose(t *testing.T) {
	fmt.Println("Test Close")
	req := requests.Requests("POST")
	for i := 0; i < 1000; i++ {
		_, err := req.Run(
			"http://localhost:1337/requests",
			requests.Datas{"SrcIp": "4312"})
		fmt.Printf("\r%d %v", i, err)
		req.Close()
	}

	spew.Dump(req)
	fmt.Println("1000 times get test end.")
}
func TestTimeout(t *testing.T) {
	println("Test Timeout")
	req := requests.Requests("GET").SetTimeout(20)
	req.SetMethod("GET").Run("http://golang.org")
}
