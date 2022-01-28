package examples

import (
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/asmcos/requests"
)

// TestBuildRequest
func TestBuildRequest(t *testing.T) {
	req, err := requests.BuildRequest("post", "http://baidu.com/a/b/c", requests.Json{
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
	req, _ := requests.BuildRequest("post", "https://baidu.com/path?q=curl&v=1", requests.Json{
		"age": 1,
	})
	curl := requests.BuildCurlRequest(req)
	if !regexp.MustCompile(`^curl -X POST .+ 'https://baidu.com/path\?q=curl&v=1'`).MatchString(curl) {
		t.Fatal(`bad curl cmd: ` + curl)
	}
	t.Log(curl)
}
