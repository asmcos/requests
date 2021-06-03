package requests

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"runtime"
)

var VERSION string = "v0.0.0"

func init() {
	_, filename, _, _ := runtime.Caller(0)
	versionFile := path.Dir(filename) + "/version"
	version, _ := ioutil.ReadFile(versionFile)
	VERSION = string(version)
}

// open file for post upload files
func openFile(filename string) *os.File {
	r, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return r
}

// handle URL params
func buildURLParams(userURL string, params ...map[string]string) (string, error) {
	parsedURL, err := url.Parse(userURL)

	if err != nil {
		return "", err
	}

	parsedQuery := parsedURL.Query()

	for _, param := range params {
		for key, value := range param {
			parsedQuery.Add(key, value)
		}
	}
    return parsedURL.String(), nil
}

