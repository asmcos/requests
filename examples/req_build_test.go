package examples

import (
	"io/ioutil"
	"testing"

	r "github.com/ahuigo/requests"
)

// TestBuildRequest
func TestBuildRequest(t *testing.T) {
	req, err := r.BuildRequest("post", "http://baidu.com/a/b/c", r.Json{
		"age": 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(req.Body)
	expectedBody := `{"age":1}`
	if string(body) != expectedBody {
		t.Fatal("Failed to build request")
	}
}

func TestBuildCurlRequest(t *testing.T) {
	// req, err := r.BuildRequest("post", "http://baidu.com/a/b/c", r.Json{
	// 	"age": 1,
	// })
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// cmd:=r.BuildCurlRequest(req)
	// if !regexp.MustCompile(``).MatchString(cmd){
	// 	t.Fatal(`bad curl cmd`)
	// }
}
