package examples

import (
	"fmt"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

func TestGetDebug(t *testing.T) {
	println("4. Get: SetDebug")
	req := requests.Requests("get").SetDebug(true)
	resp, err := req.Run("https://httpbin.org/gzip")
	if err == nil {
		fmt.Println(resp.Text())
	}
}
