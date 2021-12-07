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
	"errors"
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
var gHeader = map[string]string{
	"User-Agent": "Go-requests-" + VERSION,
}

// SetRespHandler
func SetRespHandler(fn func(*Response)) {
	respHandler = fn
}

type Session struct {
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
type Method string

// Sessions
// @params method  GET|POST|PUT|DELETE|PATCH
func Sessions() *Session {

	session := new(Session)
	session.reset()

	session.Client = &http.Client{}

	// cookiejar.New source code return jar, nil
	jar, _ := cookiejar.New(nil)

	session.Client.Jar = jar

	return session
}

// Set global header
func SetHeader(key, value string) {
	if value == "" {
		delete(gHeader, key)
		return
	}
	gHeader[key] = value
}

func (session *Session) reset() {
	session.httpreq = &http.Request{
		Method:     "GET",
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	session.Header = &session.httpreq.Header
	for key, value := range gHeader {
		session.httpreq.Header.Set(key, value)
	}
}

func (session *Session) RequestDebug() {
	if !session.debug {
		return
	}
	fmt.Println("===========Go RequestDebug !============")
	curl := BuildCurlRequest(session.httpreq)
	fmt.Println(curl)
	message, err := httputil.DumpRequestOut(session.httpreq, false)
	if err != nil {
		return
	}
	fmt.Println(string(message))

	if len(session.Client.Jar.Cookies(session.httpreq.URL)) > 0 {
		fmt.Println("Cookies:")
		for _, cookie := range session.Client.Jar.Cookies(session.httpreq.URL) {
			fmt.Println(cookie)
		}
	}
}

// cookies
// cookies only save to Client.Jar
// session.Cookies is temporary
func (session *Session) SetCookie(cookie *http.Cookie) *Session {
	session.Cookies = append(session.Cookies, cookie)
	return session
}

func (session *Session) ClearCookies() {
	session.Cookies = session.Cookies[0:0]
}

// ClientSetCookies -
func (session *Session) ClientSetCookies() {
	if len(session.Cookies) > 0 {
		// 1. Cookies have content, Copy Cookies to Client.jar
		// for _, cookie := range session.Cookies {
		// 	session.httpreq.AddCookie(cookie)
		// }
		session.Client.Jar.SetCookies(session.httpreq.URL, session.Cookies)
		// 2. Clear  Cookies
		session.ClearCookies()
	}

}

// set timeout s = second
func (session *Session) SetTimeout(n time.Duration) *Session {
	session.Client.Timeout = time.Duration(n * time.Second)
	return session
}

func (session *Session) Close() {
	session.httpreq.Close = true
}

func (session *Session) Proxy(proxyurl string) {
	urli := url.URL{}
	urlproxy, err := urli.Parse(proxyurl)
	if err != nil {
		fmt.Println("Set proxy failed")
		return
	}
	session.Client.Transport = &http.Transport{
		Proxy:           http.ProxyURL(urlproxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

// SetRespHandler -
func (session *Session) SetRespHandler(fn func(*Response)) *Session {
	session.respHandler = fn
	return session
}

// SetMethod
func (session *Session) SetMethod(method string) *Session {
	session.httpreq.Method = strings.ToUpper(method)
	return session
}

// SetHeader
func (session *Session) SetHeader(key, value string) *Session {
	session.Header.Set(key, value)
	return session
}

// BuildRequest
func (session *Session) BuildRequest(origurl string, args ...interface{}) (*http.Request, error) {
	contentType := "application/x-www-form-urlencoded"
	params := []map[string]string{}
	datas := []map[string]string{} // form data
	files := []map[string]string{} //file data
	bodyBytes := []byte{}

	for _, arg := range args {
		switch a := arg.(type) {
		// arg is Header , set to request header
		case Method:
			session.httpreq.Method = strings.ToUpper(string(a))
		case Header:
			for k, v := range a {
				session.httpreq.Header.Set(k, v)
			}
		case Params:
			params = append(params, a)
		case Datas: //Post form data,packaged in body.
			datas = append(datas, a)
		case Files:
			files = append(files, a)
		case Auth:
			session.httpreq.SetBasicAuth(a[0], a[1])
		case string:
			bodyBytes = []byte(a)
		case []byte:
			bodyBytes = a
		case *http.Cookie:
			session.SetCookie(a)
		case Json:
			contentType = "application/json"
			bodyBytes = session.buildJSON(a)
		default:
			contentType = "application/json"
			bodyBytes = session.buildJSON(a)
		}
	}
	if session.httpreq.Header.Get("Content-Type") == "" {
		session.httpreq.Header.Set("Content-Type", contentType)
	}

	disturl, _ := buildURLParams(origurl, params...)

	if len(files) > 0 {
		session.buildFilesAndForms(files, datas)
	} else if len(bodyBytes) > 0 {
		// fmt.Printf("jsonBytes=%#v\n", string(jsonBytes))
		session.setBodyBytes(bodyBytes) // set forms to body
	} else {
		Forms := session.buildForms(datas...)
		session.setBodyForms(Forms) // set forms to body
	}
	//prepare to Do
	URL, err := url.Parse(disturl)
	if err != nil {
		return nil, err
	}
	session.httpreq.URL = URL

	session.ClientSetCookies()
	// fmt.Printf("session:%#v\n", session.httpreq)
	// fmt.Printf("session-url:%#v\n", session.httpreq.URL.String())
	return session.httpreq, nil

}

// Post -
func (session *Session) Run(origurl string, args ...interface{}) (resp *Response, err error) {
	session.BuildRequest(origurl, args...)
	session.RequestDebug()
	res, err := session.Client.Do(session.httpreq)

	if err != nil {
		return nil, errors.New(session.httpreq.Method + " " + origurl + " " + err.Error())
	}

	resp = &Response{}
	resp.R = res
	req_dup := *session
	resp.session = &req_dup
	resp.ResponseDebug()
	resp.Content()
	session.reset()
	if respHandler != nil {
		respHandler(resp)
	}
	if session.respHandler != nil {
		session.respHandler(resp)
	}
	return resp, nil
}

// only set forms
func (session *Session) setBodyForms(Forms url.Values) {
	data := Forms.Encode()
	session.httpreq.Body = ioutil.NopCloser(strings.NewReader(data))
	session.httpreq.ContentLength = int64(len(data))
}

// only set forms
func (session *Session) setBodyBytes(data []byte) {
	session.httpreq.Body = ioutil.NopCloser(bytes.NewReader(data))
	session.httpreq.ContentLength = int64(len(data))
}

// upload file and form
// build to body format
func (session *Session) buildFilesAndForms(files []map[string]string, datas []map[string]string) {

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
	session.httpreq.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
	session.httpreq.ContentLength = int64(b.Len())
	session.Header.Set("Content-Type", w.FormDataContentType())
}

// build post Form data
func (session *Session) buildForms(datas ...map[string]string) (Forms url.Values) {
	Forms = url.Values{}
	for _, data := range datas {
		for key, value := range data {
			Forms.Add(key, value)
		}
	}
	return Forms
}

func (session *Session) buildJSON(data interface{}) []byte {
	jsonBytes, _ := json.Marshal(data)

	// fmt.Printf("a1=%#v,jsons=%#v\nahui\n", data, string(jsonBytes))
	return jsonBytes
}

func (session *Session) SetDebug(debug bool) *Session {
	session.debug = debug
	return session
}
