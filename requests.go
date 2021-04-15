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
	"compress/gzip"
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
	"os"
	"strings"
	"time"
)

var VERSION string = "0.8"

type Request struct {
	httpreq *http.Request
	Header  *http.Header
	Client  *http.Client
	Debug   int
	Cookies []*http.Cookie
}

type Response struct {
	R       *http.Response
	content []byte
	text    string
	req     *Request
}

type Header map[string]string
type Params map[string]string
type Datas map[string]string      // for post form
type Jsons map[string]interface{} // for Json
type AnyData interface{}          // for AnyData
type Files map[string]string      // name ,filename

// Auth - {username,password}
type Auth []string

// Requests
// @params method  GET|POST|PUT|DELETE|PATCH
func Requests(method string) *Request {

	req := new(Request)

	req.httpreq = &http.Request{
		Method:     strings.ToUpper(method),
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	req.Header = &req.httpreq.Header
	req.httpreq.Header.Set("User-Agent", "Go-Requests "+VERSION)

	req.Client = &http.Client{}

	// auto with Cookies
	// cookiejar.New source code return jar, nil
	jar, _ := cookiejar.New(nil)

	req.Client.Jar = jar

	return req
}

// Get ,req.Get

func Get(origurl string, args ...interface{}) (resp *Response, err error) {
	// call request Get
	resp, err = Requests("GET").Run(origurl, args...)
	return resp, err
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

func (req *Request) RequestDebug() {

	if req.Debug != 1 {
		return
	}

	fmt.Println("===========Go RequestDebug ============")

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
func (req *Request) SetCookie(cookie *http.Cookie) {
	req.Cookies = append(req.Cookies, cookie)
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
func (req *Request) SetTimeout(n time.Duration) {
	req.Client.Timeout = time.Duration(n * time.Second)
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

/**************/
func (resp *Response) ResponseDebug() {
	if resp.req.Debug != 1 {
		return
	}

	fmt.Println("===========Go ResponseDebug ============")

	message, err := httputil.DumpResponse(resp.R, false)
	if err != nil {
		return
	}

	fmt.Println(string(message))

}

func (resp *Response) Content() []byte {
	var err error
	if len(resp.content) > 0 {
		return resp.content
	}
	defer resp.R.Body.Close()

	var Body = resp.R.Body
	if resp.R.Header.Get("Content-Encoding") == "gzip" && resp.req.Header.Get("Accept-Encoding") != "" {
		// fmt.Println("gzip")
		reader, err := gzip.NewReader(Body)
		if err != nil {
			return nil
		}
		Body = reader
	}

	resp.content, err = ioutil.ReadAll(Body)
	if err != nil {
		return nil
	}

	return resp.content
}

func (resp *Response) Text() string {
	if resp.content == nil {
		resp.Content()
	}
	resp.text = string(resp.content)
	return resp.text
}

func (resp *Response) SaveFile(filename string) error {
	if resp.content == nil {
		resp.Content()
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(resp.content)
	f.Sync()

	return err
}

func (resp *Response) Json(v interface{}) error {
	if resp.content == nil {
		resp.Content()
	}
	return json.Unmarshal(resp.content, v)
}

func (resp *Response) Cookies() (cookies []*http.Cookie) {
	httpreq := resp.req.httpreq
	client := resp.req.Client

	cookies = client.Jar.Cookies(httpreq.URL)

	return cookies

}

/**************post*************************/
// call req.Post ,only for easy
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

// cleanup
func (req *Request) cleanup() {
	req.httpreq.Body = nil
	req.httpreq.GetBody = nil
	req.httpreq.ContentLength = 0
	//reset Cookies,
	//Client.Do can copy cookie from client.Jar to req.Header
	delete(req.httpreq.Header, "Cookie")
	// req.ClearCookies()
}

// SetMethod
func (req *Request) SetMethod(method string) *Request {
	req.httpreq.Method = strings.ToUpper(method)
	return req
}

// Post -
func (req *Request) Run(origurl string, args ...interface{}) (resp *Response, err error) {
	// cleanup
	req.cleanup()

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
				req.Header.Set(k, v)
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
			contentType = "application/text"
			bodyBytes = []byte(a)
		default:
			contentType = "application/json"
			bodyBytes = req.buildJSON(a)
		}
	}
	req.Header.Add("Content-Type", contentType)

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
		fmt.Println(err)
		return nil, err
	}

	resp = &Response{}
	resp.R = res
	resp.req = req
	resp.ResponseDebug()
	resp.Content()
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
func (req *Request) setBodyRawBytes(read io.ReadCloser) {
	req.httpreq.Body = read
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

// open file for post upload files

func openFile(filename string) *os.File {
	r, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return r
}
