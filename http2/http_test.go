package http2

import (
	"testing"
	"fmt"
	"net/url"
)

func Test_JSONPost(t *testing.T) {


	var p = url.Values{}
	p.Add("username", "smartwalle")

	var a, e = JSONRequest("POST", "http://api.smoktech.com/user/signin", p)
	fmt.Println(a, e)
}