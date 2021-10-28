package requests

import "net/http"

// BuildRequest -
func BuildRequest(method string, origurl string, args ...interface{}) (req *http.Request, err error) {
	// call request Get
	args = append(args, Method(method))
	req, err = Sessions().BuildRequest(origurl, args...)
	return
}

func BuildCurlRequest(req *http.Request) (curl string) {
	// call request Get
	return
}
