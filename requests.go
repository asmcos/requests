/* Copyright（2） 2018 by  asmcos .
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

var VERSION string = "0.9"

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
type Datas map[string]string // for post form
type Files map[string]string // name ,filename

// {username,password}
type Auth []string

func Requests() *Request {

	req := new(Request)

	req.httpreq = &http.Request{
		Method:     "GET",
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
	req.httpreq.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(req.httpreq.Body), nil
	}

	return req
}

// Get ,req.Get

func Get(origurl string, args ...interface{}) (resp *Response, err error) {
	req := Requests()

	// call request Get
	resp, err = req.Get(origurl, args...)
	return resp, err
}

func (req *Request) Get(origurl string, args ...interface{}) (resp *Response, err error) {

	req.httpreq.Method = "GET"
	req.resetBody()

	// set params ?a=b&b=c
	//set Header
	params := []map[string]string{}

	//reset Cookies,
	//Client.Do can copy cookie from client.Jar to req.Header
	delete(req.httpreq.Header, "Cookie")

	for _, arg := range args {
		switch a := arg.(type) {
		// arg is Header , set to request header
		case Header:

			req.applyHeader(a)
			// arg is "GET" params
			// ?title=website&id=1860&from=login
		case Params:
			params = append(params, a)
		case Auth:
			req.applyAuth(a)
		}
	}

	disturl, err := buildURLParams(origurl, params...)
	if err != nil {
		return nil, err
	}

	//prepare to Do
	URL, err := url.Parse(disturl)
	if err != nil {
		return nil, err
	}
	req.httpreq.URL = URL

	return req.doRequest()
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
	parsedURL.RawQuery = parsedQuery.Encode()
	return parsedURL.String(), nil
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

	urlproxy, err := url.Parse(proxyurl)
	if err != nil {
		return
	}
	t := req.ensureTransport()
	t.Proxy = http.ProxyURL(urlproxy)
}

// SetInsecureSkipVerify skips TLS certificate verification.
// Equivalent to python requests verify=False. Use only in test environments.
func (req *Request) SetInsecureSkipVerify(skip bool) {
	t := req.ensureTransport()
	if t.TLSClientConfig == nil {
		t.TLSClientConfig = &tls.Config{}
	}
	t.TLSClientConfig.InsecureSkipVerify = skip
}

func (req *Request) ensureTransport() *http.Transport {
	if req.Client.Transport == nil {
		if dt, ok := http.DefaultTransport.(*http.Transport); ok {
			t := dt.Clone()
			req.Client.Transport = t
			return t
		}
		t := &http.Transport{}
		req.Client.Transport = t
		return t
	}
	if t, ok := req.Client.Transport.(*http.Transport); ok {
		return t
	}
	t := &http.Transport{}
	req.Client.Transport = t
	return t
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

	if resp.content != nil {
		return resp.content
	}

	var Body = resp.R.Body
	if resp.R.Header.Get("Content-Encoding") == "gzip" && resp.req.Header.Get("Accept-Encoding") != "" {
		reader, err := gzip.NewReader(Body)
		if err != nil {
			resp.content = []byte{}
			return resp.content
		}
		defer reader.Close()
		Body = reader
	}

	resp.content, err = ioutil.ReadAll(Body)
	if err != nil {
		if resp.content == nil {
			resp.content = []byte{}
		}
		return resp.content
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
	req := Requests()

	// call request Get
	resp, err = req.Post(origurl, args...)
	return resp, err
}

func PostJson(origurl string, args ...interface{}) (resp *Response, err error) {
	req := Requests()

	// call request Get
	resp, err = req.PostJson(origurl, args...)
	return resp, err
}

// POST requests

func (req *Request) PostJson(origurl string, args ...interface{}) (resp *Response, err error) {

	req.httpreq.Method = "POST"
	req.resetBody()

	req.Header.Set("Content-Type", "application/json")

	//reset Cookies,
	//Client.Do can copy cookie from client.Jar to req.Header
	delete(req.httpreq.Header, "Cookie")

	params := []map[string]string{}

	for _, arg := range args {
		switch a := arg.(type) {
		// arg is Header , set to request header
		case Header:

			req.applyHeader(a)
		case Params:
			params = append(params, a)
		case string:
			req.setBody([]byte(a))
		case []byte:
			req.setBody(a)
		case Auth:
			req.applyAuth(a)
		default:
			b, err := json.Marshal(a)
			if err != nil {
				return nil, err
			}
			req.setBody(b)
		}
	}

	disturl, err := buildURLParams(origurl, params...)
	if err != nil {
		return nil, err
	}

	//prepare to Do
	URL, err := url.Parse(disturl)
	if err != nil {
		return nil, err
	}
	req.httpreq.URL = URL

	resp, err = req.doRequest()
	req.resetBody()
	return resp, err
}

func (req *Request) Post(origurl string, args ...interface{}) (resp *Response, err error) {

	req.httpreq.Method = "POST"
	req.resetBody()

	// set params ?a=b&b=c
	//set Header
	params := []map[string]string{}
	datas := []map[string]string{} // POST
	files := []map[string]string{} //post file

	//reset Cookies,
	//Client.Do can copy cookie from client.Jar to req.Header
	delete(req.httpreq.Header, "Cookie")

	for _, arg := range args {
		switch a := arg.(type) {
		// arg is Header , set to request header
		case Header:

			req.applyHeader(a)
			// arg is "GET" params
			// ?title=website&id=1860&from=login
		case Params:
			params = append(params, a)

		case Datas: //Post form data,packaged in body.
			datas = append(datas, a)
		case Files:
			files = append(files, a)
		case Auth:
			req.applyAuth(a)
		}
	}

	// default Content-Type only when caller did not set one
	if len(files) == 0 && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	disturl, err := buildURLParams(origurl, params...)
	if err != nil {
		return nil, err
	}

	if len(files) > 0 {
		err = req.buildFilesAndForms(files, datas)
		if err != nil {
			return nil, err
		}

	} else {
		Forms := req.buildForms(datas...)
		req.setBodyBytes(Forms) // set forms to body
	}
	//prepare to Do
	URL, err := url.Parse(disturl)
	if err != nil {
		return nil, err
	}
	req.httpreq.URL = URL

	resp, err = req.doRequest()
	req.resetBody()
	return resp, err
}

func (req *Request) applyHeader(h Header) {
	for k, v := range h {
		if strings.EqualFold(k, "Host") {
			req.httpreq.Host = v
			continue
		}
		req.Header.Set(k, v)
	}
}

func (req *Request) applyAuth(a Auth) {
	if len(a) >= 2 {
		req.httpreq.SetBasicAuth(a[0], a[1])
	}
}

func (req *Request) resetBody() {
	req.httpreq.Body = nil
	req.httpreq.GetBody = nil
	req.httpreq.ContentLength = 0
}

func (req *Request) doRequest() (*Response, error) {
	req.ClientSetCookies()
	req.RequestDebug()

	httpReq := req.httpreq
	if req.httpreq.GetBody != nil {
		cloned := req.httpreq.Clone(req.httpreq.Context())
		body, err := req.httpreq.GetBody()
		if err != nil {
			return nil, err
		}
		cloned.Body = body
		httpReq = cloned
	}

	res, err := req.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	resp := &Response{}
	resp.R = res
	resp.req = req

	resp.Content()
	res.Body.Close()

	resp.ResponseDebug()
	return resp, nil
}

func (req *Request) setBody(data []byte) {
	req.httpreq.ContentLength = int64(len(data))
	req.httpreq.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(bytes.NewReader(data)), nil
	}
	req.httpreq.Body, _ = req.httpreq.GetBody()
}

// only set forms
func (req *Request) setBodyBytes(Forms url.Values) {

	// maybe
	data := Forms.Encode()
	req.setBody([]byte(data))
}

// upload file and form
// build to body format
func (req *Request) buildFilesAndForms(files []map[string]string, datas []map[string]string) error {

	//handle file multipart

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// write form fields before files so S3-style uploads receive required keys first
	for _, data := range datas {
		for k, v := range data {
			if err := w.WriteField(k, v); err != nil {
				return err
			}
		}
	}

	for _, file := range files {
		for k, v := range file {
			part, err := w.CreateFormFile(k, v)
			if err != nil {
				return err
			}
			f, err := os.Open(v)
			if err != nil {
				return err
			}
			_, err = io.Copy(part, f)
			f.Close()
			if err != nil {
				return err
			}
		}
	}

	if err := w.Close(); err != nil {
		return err
	}
	// set file header example:
	// "Content-Type": "multipart/form-data; boundary=------------------------7d87eceb5520850c",
	req.setBody(b.Bytes())
	req.Header.Set("Content-Type", w.FormDataContentType())
	return nil
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
