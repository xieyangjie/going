package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"io"
)

const (
	K_HTTP_METHOD_GET    = "GET"
	K_HTTP_METHOD_POST   = "POST"
	K_HTTP_METHOD_HEAD   = "HEAD"
	K_HTTP_METHOD_PUT    = "PUT"
	K_HTTP_METHOD_DELETE = "DELETE"
)

////////////////////////////////////////////////////////////////////////////////
type Client struct {
	//请求方法
	method string
	//url
	urlString string
	//timeout
	timeout time.Duration

	//参数
	params url.Values
	//请求头
	headers map[string]string

	body string

	username string
	password string
}

func NewClient() *Client {
	var client = &Client{}
	client.method = K_HTTP_METHOD_GET
	return client
}

func (this *Client) SetMethod(method string) {
	this.method = strings.ToUpper(method)
}

func (this *Client) SetURLString(urlString string) {
	this.urlString = urlString
}

func (this *Client) SetTimeout(timeout time.Duration) {
	this.timeout = timeout
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

func (this *Client) SetBody(body string) {
	this.body = body
}

func (this *Client) SetBasicAuth(username, password string) {
	this.username = username
	this.password = password
}

func (this *Client) createGetURL() string {
	var paramStr = "?"
	if this.params != nil {
		paramStr = paramStr + this.params.Encode()
	}
	return this.urlString + paramStr
}

func (this *Client) createGetRequest() (*http.Request, error) {
	return http.NewRequest(this.method, this.createGetURL(), nil)
}

func (this *Client) createPostRequest() (*http.Request, error) {
	if _, ok := this.headers["Content-Type"]; !ok {
		this.headers["Content-Type"] = "application/x-www-form-urlencoded"
	}

	var body io.Reader
	if len(this.body) == 0 {
		body = strings.NewReader(this.params.Encode())
	} else {
		body = strings.NewReader(this.body)
	}

	return http.NewRequest(this.method, this.urlString, body)
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

	if len(this.username) > 0 {
		request.SetBasicAuth(this.username, this.password)
	}

	for key, value := range this.headers {
		request.Header.Set(key, value)
	}

	var c = &http.Client{}
	c.Timeout = this.timeout

	return c.Do(request)
}

func (this *Client) DoRequest() ([]byte, error) {
	responseHandler, err := this.doRequest()
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

func (this *Client) DoJsonRequest() (result map[string]interface{}, err error) {
	err = this.JsonRequest(&result)
	return result, err
}

func (this *Client) JsonRequest(result interface{}) (error) {
	respBodyByte, err := this.DoRequest()
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBodyByte, &result)
	return err
}

////////////////////////////////////////////////////////////////////////////////
func DoRequest(method string, urlString string, param map[string]string) ([]byte, error) {
	var client = NewClient()
	client.SetMethod(method)
	client.SetURLString(urlString)
	for key, value := range param {
		client.SetParam(key, value)
	}
	return client.DoRequest()
}

func DoGet(urlString string, param map[string]string) ([]byte, error) {
	return DoRequest(K_HTTP_METHOD_GET, urlString, param)
}

func DoPost(urlString string, param map[string]string) ([]byte, error) {
	return DoRequest(K_HTTP_METHOD_POST, urlString, param)
}

func DoJSONRequest(method string, urlString string, param map[string]string) (map[string]interface{}, error) {
	var client = NewClient()
	client.SetMethod(method)
	client.SetURLString(urlString)
	for key, value := range param {
		client.SetParam(key, value)
	}
	return client.DoJsonRequest()
}

func DoJSONGet(urlString string, param map[string]string) (map[string]interface{}, error) {
	return DoJSONRequest(K_HTTP_METHOD_GET, urlString, param)
}

func DoJSONPost(urlString string, param map[string]string) (map[string]interface{}, error) {
	return DoJSONRequest(K_HTTP_METHOD_POST, urlString, param)
}
