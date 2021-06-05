package examples

import (
	"testing"

	r "github.com/ahuigo/requests"
)

// Test Session with cookies
func TestSessionWithCookie(t *testing.T) {
	var data struct {
		Cookies struct {
			Count string `json:"count"`
		}
	}
	session := r.Sessions()
	// set cookies: count=100
	session.Get("https://httpbin.org/cookies/set?count=100")
	// get cookies
	resp, err := session.Get("https://httpbin.org/cookies")
	if err == nil {
		resp.Json(&data)
		if data.Cookies.Count != "100" {
			t.Fatal("Failed to get valid cookies: " + resp.Text())
		}
	}
	if err != nil {
		t.Fatal(err)
	}
}
