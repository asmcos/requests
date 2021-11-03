package requests

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/alessio/shellescape"
)

// BuildRequest -
func BuildRequest(method string, origurl string, args ...interface{}) (req *http.Request, err error) {
	// call request Get
	args = append(args, Method(method))
	req, err = Sessions().BuildRequest(origurl, args...)
	return
}

func BuildCurlRequest(req *http.Request) (curl string) {
	curl = "curl -X " + req.Method + " "
	// req.Host + req.URL.Path + "?" + req.URL.RawQuery + " " + req.Proto + " "
	headers := getHeaders(req)
	for _, kv := range *headers {
		curl += `-H ` + shellescape.Quote(kv[0]+": "+kv[1]) + ` `
	}
	// // cookies
	// for _, cookie := range req.Cookies() {
	// 	fmt.Printf("cookie:%#v\n", cookie)
	// }
	// body
	buf, _ := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf)) // important!!
	if len(buf) > 0 {
		curl += `-d ` + shellescape.Quote(string(buf))
	}

	curl += " " + shellescape.Quote(req.URL.String())
	return curl
}

// getHeaders
func getHeaders(req *http.Request) *[][2]string {
	headers := [][2]string{}
	for k, vs := range req.Header {
		for _, v := range vs {
			headers = append(headers, [2]string{k, v})
		}
	}
	n := len(headers)
	// fmt.Printf("%#v\n", headers)
	// sort headers
	for i := 0; i < n; i++ {
		for j := n - 1; j > i; j-- {
			jj := j - 1
			h1, h2 := headers[j], headers[jj]
			if h1[0] < h2[0] {
				headers[jj], headers[j] = headers[j], headers[jj]
			}
		}
	}
	return &headers
}
