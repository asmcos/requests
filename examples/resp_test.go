package examples

import (
	"testing"

	"github.com/asmcos/requests"
	_ "github.com/asmcos/requests/init"
)

// Test response headers
func TestResponseHeader(t *testing.T) {
	resp, _ := requests.Get("https://httpbin.org/get")
	println("content-type:", resp.R.Header.Get("content-type"))
	//println(resp.Text())
}

// Test response body
func TestResponseBody(t *testing.T) {
	resp, _ := requests.Get("https://httpbin.org/get")
	println(resp.Text())
}
