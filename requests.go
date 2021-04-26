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

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var respHandler func(*Response)

// SetRespHandler
func SetRespHandler(fn func(*Response)) {
	respHandler = fn
}

type Request struct {
	httpreq     *http.Request
	Client      *http.Client
	debug       bool
	respHandler func(*Response)
	// global header
	Header  *http.Header
	Cookies []*http.Cookie
}

type Header map[string]string
type Params map[string]string
type Datas map[string]string     // for post form
type Json map[string]interface{} // for Json
type Files map[string]string     // name ,filename
// type AnyData interface{}         // for AnyData

// Auth - {username,password}
type Auth []string

// Requests
// @params method  GET|POST|PUT|DELETE|PATCH
func Requests() *Request {

	req := new(Request)
	req.reset()

	req.Client = &http.Client{}

	// cookiejar.New source code return jar, nil
	jar, _ := cookiejar.New(nil)

	req.Client.Jar = jar

	return req
}

func (req *Request) reset() {
	req.httpreq = &http.Request{
		Method:     "GET",
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	req.Header = &req.httpreq.Header
	req.httpreq.Header.Set("User-Agent", "Go-Requests "+VERSION)
}

func (req *Request) RequestDebug() {
	if !req.debug {
		return
	}
	fmt.Println("===========Go RequestDebug !============")

	message, err := httputil.DumpRequestOut(req.httpreq, false)
	if err != nil {
		return
	}
	fmt.Println(string(message))

	if len(req.Client.Jar.Cookies(req.httpreq.URL)) > 0 {
		fmt.Println("Cookies:")
		for _, cookie := range req.Client.Jar.Cookies(req.httpreq.URL) {
			fmt.Println(cookie)
		}
	}
}

// cookies
// cookies only save to Client.Jar
// req.Cookies is temporary
func (req *Request) SetCookie(cookie *http.Cookie) *Request {
	req.Cookies = append(req.Cookies, cookie)
	return req
}

func (req *Request) ClearCookies() {
	req.Cookies = req.Cookies[0:0]
}

// ClientSetCookies -
func (req *Request) ClientSetCookies() {
	if len(req.Cookies) > 0 {
		// 1. Cookies have content, Copy Cookies to Client.jar
		// 2. Clear  Cookies
		req.Client.Jar.SetCookies(req.httpreq.URL, req.Cookies)
		req.ClearCookies()
	}

}

// set timeout s = second
func (req *Request) SetTimeout(n time.Duration) *Request {
	req.Client.Timeout = time.Duration(n * time.Second)
	return req
}

func (req *Request) Close() {
	req.httpreq.Close = true
}

func (req *Request) Proxy(proxyurl string) {
	urli := url.URL{}
	urlproxy, err := urli.Parse(proxyurl)
	if err != nil {
		fmt.Println("Set proxy failed")
		return
	}
	req.Client.Transport = &http.Transport{
		Proxy:           http.ProxyURL(urlproxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

// cleanup delete
func (req *Request) cleanup() {
	req.httpreq.Body = nil
	req.httpreq.GetBody = nil
	req.httpreq.ContentLength = 0
	// req.Header = &http.Header{}
	// req.Header = &req.httpreq.Header

	//reset Cookies,
	//Client.Do can copy cookie from client.Jar to req.Header
	delete(req.httpreq.Header, "Cookie")
	// req.ClearCookies()
}

// SetRespHandler
func (req *Request) SetRespHandler(fn func(*Response)) *Request {
	req.respHandler = fn
	return req
}

// SetMethod
func (req *Request) SetMethod(method string) *Request {
	req.httpreq.Method = strings.ToUpper(method)
	return req
}

// SetHeader
func (req *Request) SetHeader(key, value string) *Request {
	req.Header.Set(key, value)
	return req
}

// Post -
func (req *Request) Run(origurl string, args ...interface{}) (resp *Response, err error) {
	// cleanup
	// req.cleanup()

	// set params ?a=b&b=c
	//set Header
	contentType := "application/x-www-form-urlencoded"
	params := []map[string]string{}
	datas := []map[string]string{} // form data
	files := []map[string]string{} //file data
	bodyBytes := []byte{}

	for _, arg := range args {
		switch a := arg.(type) {
		// arg is Header , set to request header
		case Header:
			for k, v := range a {
				req.httpreq.Header.Set(k, v)
			}
		case Params:
			params = append(params, a)
		case Datas: //Post form data,packaged in body.
			datas = append(datas, a)
		case Files:
			files = append(files, a)
		case Auth:
			req.httpreq.SetBasicAuth(a[0], a[1])
		case string:
			bodyBytes = []byte(a)
		case Json:
			contentType = "application/json"
			bodyBytes = req.buildJSON(a)
		default:
			contentType = "application/json"
			bodyBytes = req.buildJSON(a)
		}
	}
	if req.httpreq.Header.Get("Content-Type") == "" {
		req.httpreq.Header.Set("Content-Type", contentType)
	}

	disturl, _ := buildURLParams(origurl, params...)

	if len(files) > 0 {
		req.buildFilesAndForms(files, datas)
	} else if len(bodyBytes) > 0 {
		// fmt.Printf("jsonBytes=%#v\n", string(jsonBytes))
		req.setBodyBytes(bodyBytes) // set forms to body
	} else {
		Forms := req.buildForms(datas...)
		req.setBodyForms(Forms) // set forms to body
	}
	//prepare to Do
	URL, err := url.Parse(disturl)
	if err != nil {
		return nil, err
	}
	req.httpreq.URL = URL

	req.ClientSetCookies()

	req.RequestDebug()

	// fmt.Printf("req:%#v\n", req.httpreq)
	// fmt.Printf("req-url:%#v\n", req.httpreq.URL.String())
	res, err := req.Client.Do(req.httpreq)

	if err != nil {
		return nil, err
	}

	resp = &Response{}
	resp.R = res
	req_dup := *req
	resp.req = &req_dup
	resp.ResponseDebug()
	resp.Content()
	req.reset()
	if respHandler != nil {
		respHandler(resp)
	}
	if req.respHandler != nil {
		req.respHandler(resp)
	}
	return resp, nil
}

// only set forms
func (req *Request) setBodyForms(Forms url.Values) {
	data := Forms.Encode()
	req.httpreq.Body = ioutil.NopCloser(strings.NewReader(data))
	req.httpreq.ContentLength = int64(len(data))
}

// only set forms
func (req *Request) setBodyBytes(data []byte) {
	req.httpreq.Body = ioutil.NopCloser(bytes.NewReader(data))
	req.httpreq.ContentLength = int64(len(data))
}

// upload file and form
// build to body format
func (req *Request) buildFilesAndForms(files []map[string]string, datas []map[string]string) {

	//handle file multipart

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for _, file := range files {
		for k, v := range file {
			part, err := w.CreateFormFile(k, v)
			if err != nil {
				fmt.Printf("Upload %s failed!", v)
				panic(err)
			}
			file := openFile(v)
			_, err = io.Copy(part, file)
			if err != nil {
				panic(err)
			}
		}
	}

	for _, data := range datas {
		for k, v := range data {
			w.WriteField(k, v)
		}
	}

	w.Close()
	// set file header example:
	// "Content-Type": "multipart/form-data; boundary=------------------------7d87eceb5520850c",
	req.httpreq.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
	req.httpreq.ContentLength = int64(b.Len())
	req.Header.Set("Content-Type", w.FormDataContentType())
}

// build post Form data
func (req *Request) buildForms(datas ...map[string]string) (Forms url.Values) {
	Forms = url.Values{}
	for _, data := range datas {
		for key, value := range data {
			Forms.Add(key, value)
		}
	}
	return Forms
}

func (req *Request) buildJSON(data interface{}) []byte {
	jsonBytes, _ := json.Marshal(data)

	// fmt.Printf("a1=%#v,jsons=%#v\nahui\n", data, string(jsonBytes))
	return jsonBytes
}

func (req *Request) SetDebug(debug bool) *Request {
	req.debug = debug
	return req
}
