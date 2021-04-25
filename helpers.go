package requests

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
)

var VERSION string = "v0.0.0"

func init() {
	_, filename, _, _ := runtime.Caller(0)
	println("filename:" + filename)
	versionFile := path.Dir(filename) + "/version"
	version, _ := ioutil.ReadFile(versionFile)
	println("version:" + string(version))
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

	parsedQuery, err := url.ParseQuery(parsedURL.RawQuery)

	if err != nil {
		return "", nil
	}

	for _, param := range params {
		for key, value := range param {
			parsedQuery.Add(key, value)
		}
	}
	return addQueryParams(parsedURL, parsedQuery), nil
}

func addQueryParams(parsedURL *url.URL, parsedQuery url.Values) string {
	if len(parsedQuery) > 0 {
		return strings.Join([]string{strings.Replace(parsedURL.String(), "?"+parsedURL.RawQuery, "", -1), parsedQuery.Encode()}, "?")
	}
	return strings.Replace(parsedURL.String(), "?"+parsedURL.RawQuery, "", -1)
}
