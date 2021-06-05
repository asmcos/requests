/* Copyright（2） 2018 by  asmcos and ahuigo .
Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package requests

import "net/http"

// BuildRequest -
func BuildRequest(method string, origurl string, args ...interface{}) (req *http.Request, err error) {
	// call request Get
	args = append(args, Method(method))
	req, err = Sessions().BuildRequest(origurl, args...)
	return
}

/**************post/get/delete/patch*************************/
func Get(origurl string, args ...interface{}) (resp *Response, err error) {
	// call request Get
	resp, err = Sessions().Get(origurl, args...)
	return resp, err
}

func Post(origurl string, args ...interface{}) (resp *Response, err error) {
	resp, err = Sessions().Post(origurl, args...)
	return
}

// Put
func Put(origurl string, args ...interface{}) (resp *Response, err error) {
	resp, err = Sessions().Put(origurl, args...)
	return
}

// Delete
func Delete(origurl string, args ...interface{}) (resp *Response, err error) {
	resp, err = Sessions().Delete(origurl, args...)
	return
}

// Patch
func Patch(origurl string, args ...interface{}) (resp *Response, err error) {
	resp, err = Sessions().Patch(origurl, args...)
	return
}

func (session *Session) Get(origurl string, args ...interface{}) (resp *Response, err error) {
	session.httpreq.Method = "GET"
	resp, err = session.Run(origurl, args...)
	return resp, err
}
func (session *Session) Post(origurl string, args ...interface{}) (resp *Response, err error) {
	session.httpreq.Method = "POST"
	resp, err = session.Run(origurl, args...)
	return resp, err
}
func (session *Session) Delete(origurl string, args ...interface{}) (resp *Response, err error) {
	session.httpreq.Method = "DELETE"
	resp, err = session.Run(origurl, args...)
	return resp, err
}
func (session *Session) Put(origurl string, args ...interface{}) (resp *Response, err error) {
	session.httpreq.Method = "PUT"
	resp, err = session.Run(origurl, args...)
	return resp, err
}

func (session *Session) Patch(origurl string, args ...interface{}) (resp *Response, err error) {
	session.httpreq.Method = "PATCH"
	resp, err = session.Run(origurl, args...)
	return resp, err
}
