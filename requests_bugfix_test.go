package requests

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPostJsonBodyHasNoExtraBytes(t *testing.T) {
	jsonStr := `{"name":"requests_post_test"}`
	var gotBody []byte
	var contentLength int64
	var transferEncoding []string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentLength = r.ContentLength
		transferEncoding = append([]string{}, r.TransferEncoding...)
		gotBody, _ = ioutil.ReadAll(r.Body)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	resp, err := PostJson(ts.URL, jsonStr)
	if err != nil {
		t.Fatalf("PostJson error: %v", err)
	}
	if resp.Text() != `{"ok":true}` {
		t.Fatalf("unexpected response text: %q", resp.Text())
	}
	if string(gotBody) != jsonStr {
		t.Fatalf("request body mismatch\n got: %q\nwant: %q", gotBody, jsonStr)
	}
	if contentLength != int64(len(jsonStr)) {
		t.Fatalf("Content-Length = %d, want %d", contentLength, len(jsonStr))
	}
	if len(transferEncoding) > 0 {
		t.Fatalf("unexpected Transfer-Encoding %v, body should use Content-Length", transferEncoding)
	}
}

func TestPostJsonMarshalHasNoTrailingNewline(t *testing.T) {
	payload := map[string]string{"name": "requests_post_test"}
	var gotBody []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = ioutil.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_, err := PostJson(ts.URL, payload)
	if err != nil {
		t.Fatalf("PostJson error: %v", err)
	}

	want, _ := json.Marshal(payload)
	if string(gotBody) != string(want) {
		t.Fatalf("request body mismatch\n got: %q\nwant: %q", gotBody, want)
	}
	if strings.HasSuffix(string(gotBody), "\n") {
		t.Fatalf("json body should not have trailing newline: %q", gotBody)
	}
}

func TestPostJsonGetBodyCanReplay(t *testing.T) {
	req := Requests()
	body := `{"id":1}`
	req.setBody([]byte(body))

	if req.httpreq.GetBody == nil {
		t.Fatal("GetBody is nil, http2 retry / redirect cannot rewind body")
	}

	first, err := ioutil.ReadAll(req.httpreq.Body)
	if err != nil {
		t.Fatal(err)
	}
	replay, err := req.httpreq.GetBody()
	if err != nil {
		t.Fatal(err)
	}
	second, err := ioutil.ReadAll(replay)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != body || string(second) != body {
		t.Fatalf("GetBody replay failed first=%q second=%q", first, second)
	}
}

func TestPostThenGetDoesNotSendStaleBody(t *testing.T) {
	var methods []string
	var lengths []int64
	var bodies []string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		methods = append(methods, r.Method)
		lengths = append(lengths, r.ContentLength)
		bodies = append(bodies, string(b))
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	client := Requests()
	if _, err := client.Post(ts.URL, Datas{"abc": "123", "ddd": "789"}); err != nil {
		t.Fatalf("Post error: %v", err)
	}
	if _, err := client.Get(ts.URL); err != nil {
		t.Fatalf("Get error: %v", err)
	}

	if len(methods) != 2 || methods[0] != "POST" || methods[1] != "GET" {
		t.Fatalf("unexpected methods: %v", methods)
	}
	if lengths[1] != 0 {
		t.Fatalf("GET leftover ContentLength=%d body=%q", lengths[1], bodies[1])
	}
	if bodies[1] != "" {
		t.Fatalf("GET should not send POST body, got %q", bodies[1])
	}
}

func TestPostDefaultAndCustomContentType(t *testing.T) {
	var gotType string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	if _, err := Post(ts.URL, Datas{"a": "b"}); err != nil {
		t.Fatal(err)
	}
	if gotType != "application/x-www-form-urlencoded" {
		t.Fatalf("default Content-Type = %q", gotType)
	}

	req := Requests()
	req.Header.Set("Content-Type", "text/plain")
	if _, err := req.Post(ts.URL, Datas{"a": "b"}); err != nil {
		t.Fatal(err)
	}
	if gotType != "text/plain" {
		t.Fatalf("Header.Set Content-Type was overwritten: %q", gotType)
	}

	if _, err := Post(ts.URL, Header{"Content-Type": "application/custom"}, Datas{"a": "b"}); err != nil {
		t.Fatal(err)
	}
	if gotType != "application/custom" {
		t.Fatalf("Header arg Content-Type was ignored: %q", gotType)
	}
}

func TestHostHeaderSetsRequestHost(t *testing.T) {
	var gotHost string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHost = r.Host
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	if _, err := Get(ts.URL, Header{"Host": "example.test"}); err != nil {
		t.Fatal(err)
	}
	if gotHost != "example.test" {
		t.Fatalf("Host = %q, want example.test", gotHost)
	}
}

func TestBuildURLParamsPreservesAndAddsQuery(t *testing.T) {
	got, err := buildURLParams("http://example.com/path?a=1", map[string]string{"b": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "a=1") || !strings.Contains(got, "b=2") {
		t.Fatalf("unexpected url: %s", got)
	}
}

func TestMultipartWritesFieldsBeforeFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "requests-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "hello.txt")
	if err := ioutil.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	var raw []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ = ioutil.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_, err = Post(ts.URL, Datas{"key": "my-object", "acl": "private"}, Files{"file": path})
	if err != nil {
		t.Fatal(err)
	}

	body := string(raw)
	keyPos := strings.Index(body, `name="key"`)
	filePos := strings.Index(body, `name="file"`)
	if keyPos < 0 || filePos < 0 {
		t.Fatalf("multipart body missing fields:\n%s", body)
	}
	if keyPos > filePos {
		t.Fatalf("form fields must appear before files for S3-style uploads:\n%s", body)
	}
}

func TestSetInsecureSkipVerify(t *testing.T) {
	req := Requests()
	req.SetInsecureSkipVerify(true)
	tr, ok := req.Client.Transport.(*http.Transport)
	if !ok || tr.TLSClientConfig == nil || !tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("SetInsecureSkipVerify did not enable tls.Config.InsecureSkipVerify")
	}

	req.SetInsecureSkipVerify(false)
	if tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("SetInsecureSkipVerify(false) did not disable skip verify")
	}

	req.Proxy("http://127.0.0.1:8080")
	tr = req.Client.Transport.(*http.Transport)
	if tr.TLSClientConfig != nil && tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("Proxy() should not force InsecureSkipVerify")
	}
}

func TestGetConnectionErrorIsReturned(t *testing.T) {
	req := Requests()
	_, err := req.Get("http://127.0.0.1:1")
	if err == nil {
		t.Fatal("expected connection error")
	}
}

func TestGzipReaderClosesAndReads(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		gz.Write([]byte("hello gzip"))
		gz.Close()
	}))
	defer ts.Close()

	req := Requests()
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := req.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text() != "hello gzip" {
		t.Fatalf("gzip body = %q", resp.Text())
	}
}

func TestPostJsonParams(t *testing.T) {
	var gotURL string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_, err := PostJson(ts.URL, Params{"q": "1"}, `{"a":1}`)
	if err != nil {
		t.Fatal(err)
	}
	if gotURL != "q=1" {
		t.Fatalf("query = %q, want q=1", gotURL)
	}
}
