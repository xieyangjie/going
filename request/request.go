package request

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"fmt"
)

func NewRequest(method, url string, params url.Values) (*http.Request, error) {
	var m = strings.ToUpper(method)
	var body io.Reader
	if m == "GET" || m == "HEAD" {
		if strings.Contains(url, "?") == false  && len(params) > 0 {
			url = url + "?" + params.Encode()
		}
	} else {
		body = strings.NewReader(params.Encode())
	}
	return http.NewRequest(m, url, body)
}

func DoRequest(c *http.Client, req *http.Request) (*http.Response, []byte, error) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	rep, err := c.Do(req)
	if err != nil {
		return rep, nil, err
	}
	defer rep.Body.Close()

	repBody, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return rep, nil, err
	}
	return rep, repBody, err
}

func Request(method, url string, params url.Values) ([]byte, error) {
	req, err := NewRequest(method, url, params)
	if err != nil {
		return nil, err
	}
	_, repBody, err := DoRequest(http.DefaultClient, req)
	return repBody, err
}

func DoJSONRequest(c *http.Client, req *http.Request, result interface{}) (*http.Response, error) {
	rep, repBody, err := DoRequest(c, req)
	if err != nil {
		return rep, err
	}
	fmt.Println(string(repBody))
	err = json.Unmarshal(repBody, &result)
	return rep, err
}

func JSONRequest(method, url string, params url.Values) (result map[string]interface{}, err error) {
	req, err := NewRequest(method, url, params)
	if err != nil {
		return nil, err
	}
	_, err = DoJSONRequest(http.DefaultClient, req, &result)
	return result, err
}
