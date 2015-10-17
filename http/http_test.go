package http

import (
	"testing"
	"fmt"
)

func Test_JSONPost(t *testing.T) {
	var a, e = DoJSONPost("http://127.0.0.1:9010/user/signin", map[string]string{"username": "smartwalle", "password":"123456"})
	fmt.Println(a, e)
}