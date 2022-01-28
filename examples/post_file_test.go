package examples

import (
	"os"
	"testing"

	"github.com/asmcos/requests"
	_ "github.com/asmcos/requests/init"
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
    if err !=nil {
        t.Error(err)
	}
	var data = struct {
		Files struct{
            File2 string
        }
	}{}
	err = resp.Json(&data)
	if data.Files.File2== ""{
        t.Error("invalid response body:", resp.Text())
	}

}
