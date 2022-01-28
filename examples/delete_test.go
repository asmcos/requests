package examples

import (
	"fmt"
	"testing"

	"github.com/asmcos/requests"
	_ "github.com/asmcos/requests/init"
)

// Delete Form Request
func TestDeleteForm(t *testing.T) {
	println("Test DELETE method: form data(x-wwww-form-urlencoded)")
	data := requests.Datas{
		"comments": "ew",
	}
	session := requests.Sessions() //.SetDebug(true)
	resp, err := session.Delete("https://www.httpbin.org/delete", data)
	if err == nil {
		fmt.Println(resp.Text())
	}
}
