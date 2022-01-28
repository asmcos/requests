package examples

import (
	"fmt"
	"testing"

	"github.com/asmcos/requests"
	_ "github.com/asmcos/requests/init"
)

func TestAuth(t *testing.T) {
	println("3. Get: Set Auth")
	// test authentication usernae,password
	//documentation https://www.httpwatch.com/httpgallery/authentication/#showExample10
	resp, err := requests.Get(
		"https://www.httpwatch.com/httpgallery/authentication/authenticatedimage/default.aspx?0.45874470316137206",
		requests.Auth{"httpwatch", "foo"},
	)
	if err == nil {
		fmt.Println(resp.R)
	}
	// this save file test PASS
	// resp.SaveFile("auth.jpeg")
}
