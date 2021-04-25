package examples

import (
	"fmt"
	"os"
	"testing"

	"github.com/ahuigo/requests"
	_ "github.com/ahuigo/requests/init"
)

func TestPostFile(t *testing.T) {
	path, _ := os.Getwd()

	resp, err := requests.Post(
		"https://www.httpbin.org/post",
		requests.Files{
			"file1": path + "/README.md",
			"file2": path + "/version",
		},
	)
	if err == nil {
		fmt.Println(resp.Text())
	}
}
