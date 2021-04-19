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

/**************post/get/delete/patch*************************/
func Get(origurl string, args ...interface{}) (resp *Response, err error) {
	// call request Get
	resp, err = Requests("GET").Run(origurl, args...)
	return resp, err
}

func Post(origurl string, args ...interface{}) (resp *Response, err error) {
	resp, err = Requests("POST").Run(origurl, args...)
	return
}

// Put
func Put(origurl string, args ...interface{}) (resp *Response, err error) {
	resp, err = Requests("PUT").Run(origurl, args...)
	return
}

// Delete
func Delete(origurl string, args ...interface{}) (resp *Response, err error) {
	resp, err = Requests("DELETE").Run(origurl, args...)
	return
}

// Patch
func Patch(origurl string, args ...interface{}) (resp *Response, err error) {
	resp, err = Requests("PATCH").Run(origurl, args...)
	return
}

func (req *Request) Get(origurl string, args ...interface{}) (resp *Response, err error) {
	req.httpreq.Method = "GET"
	resp, err = req.Run(origurl, args...)
	return resp, err
}
func (req *Request) Post(origurl string, args ...interface{}) (resp *Response, err error) {
	req.httpreq.Method = "POST"
	resp, err = req.Run(origurl, args...)
	return resp, err
}
func (req *Request) Delete(origurl string, args ...interface{}) (resp *Response, err error) {
	req.httpreq.Method = "DELETE"
	resp, err = req.Run(origurl, args...)
	return resp, err
}
func (req *Request) Patch(origurl string, args ...interface{}) (resp *Response, err error) {
	req.httpreq.Method = "PATCH"
	resp, err = req.Run(origurl, args...)
	return resp, err
}
