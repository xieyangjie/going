package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	K_HTTP_METHOD_GET 		= "GET"
	K_HTTP_METHOD_POST 		= "POST"
	K_HTTP_METHOD_HEAD 		= "HEAD"
	K_HTTP_METHOD_PUT 		= "PUT"
	K_HTTP_METHOD_DELETE	= "DELETE"
)

func NewClient() *Client {
	var client = &Client{}
	client.method = K_HTTP_METHOD_GET
	return client
}

type Client struct {
	//请求方法
	method 		string
	//url
	urlString 	string

	//参数
	params	url.Values
	//请求头
	headers	map[string]string
}

func (this *Client) SetMethod(method string) {
	this.method = strings.ToUpper(method)
}

func (this *Client) SetURLString(urlString string) {
	this.urlString = urlString
}

func (this *Client) SetHeader(key string, value string) {
	if this.headers == nil {
		this.headers = make(map[string]string)
	}
	this.headers[key] = value
}

func (this *Client) SetParam(key string, value string) {
	if this.params == nil {
		this.params = make(url.Values)
	}
	this.params.Set(key, value)
}

func (this *Client) createGetURL() string {
	var paramStr = "?"
	if this.params != nil {
		paramStr = paramStr + this.params.Encode()
	}
	return this.urlString+paramStr
}

func (this *Client) createGetRequest() (*http.Request, error) {
	return http.NewRequest(this.method, this.createGetURL(), nil)
}

func (this *Client) createPostRequest() (*http.Request, error) {
	if _, ok := this.headers["Content-Type"]; !ok {
		this.headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
	return http.NewRequest(this.method, this.urlString, strings.NewReader(this.params.Encode()))
}

func (this *Client) doRequest() (*http.Response, error) {
	if this.headers == nil {
		this.headers = make(map[string]string)
	}

	var method = strings.ToUpper(this.method)

	var request *http.Request
	var err error

	if method == K_HTTP_METHOD_GET || method == K_HTTP_METHOD_HEAD {
		request, err = this.createGetRequest()
	} else {
		request, err = this.createPostRequest()
	}

	if err != nil {
		return nil, err
	}

	for key, value := range this.headers {
		request.Header.Set(key, value)
	}

	return http.DefaultClient.Do(request)
}

func (this *Client) DoRequest() ([]byte, error) {
	responseHandler,err := this.doRequest()
	if err != nil {
		return nil, err
	}

	defer responseHandler.Body.Close()

	respBodyByte, err := ioutil.ReadAll(responseHandler.Body)
	if err != nil {
		return nil, err
	}
	return respBodyByte, err
}

func (this *Client) DoJsonRequest() (map[string]interface{}, error) {
	var result map[string]interface{}
	respBodyByte, err := this.DoRequest()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBodyByte, &result)
	return result, err
}


func DoRequest(method string, urlString string, param map[string]string) ([]byte, error) {
	var client = NewClient()
	client.SetMethod(method)
	client.SetURLString(urlString)
	for key, value := range param {
		client.SetParam(key, value)
	}
	return client.DoRequest()
}

func DoJsonRequest(method string, urlString string, param map[string]string) (map[string]interface{}, error) {
	var client = NewClient()
	client.SetMethod(method)
	client.SetURLString(urlString)
	for key, value := range param {
		client.SetParam(key, value)
	}
	return client.DoJsonRequest()
}

func DoGet(urlString string, param map[string]string) (map[string]interface{}, error) {
	return DoJsonRequest(K_HTTP_METHOD_GET, urlString, param)
}

func DoPost(urlString string, param map[string]string) (map[string]interface{}, error) {
	return DoJsonRequest(K_HTTP_METHOD_POST, urlString, param)
}
