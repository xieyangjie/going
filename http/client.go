package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	K_HTTP_REQ_METHOD_GET  		= "GET"
	K_HTTP_REQ_METHOD_POST 		= "POST"
	K_HTTP_REQ_METHOD_HEAD 		= "HEAD"
	K_HTTP_REQ_METHOD_PUT		= "PUT"
	K_HTTP_REQ_METHOD_DELETE	= "DELETE"
)

func NewClient() *Client {
	var client = &Client{}
	client.Method = K_HTTP_REQ_METHOD_GET
	return client
}

type Client struct {
	//请求方法
	Method 		string
	//url
	URLString	string

	//参数
	params	map[string]string
	//请求头
	headers	map[string]string
}

func (this *Client) SetHeader(key string, value string) {
	if this.headers == nil {
		this.headers = make(map[string]string)
	}
	this.headers[key] = value
}

func (this *Client) SetParam(key string, value string) {
	if this.params == nil {
		this.params = make(map[string]string)
	}
	this.params[key] = value
}

func (this *Client) createGetURL() string {
	var paramStr = "?"
	for key, value := range this.params {
		paramStr = paramStr + url.QueryEscape(key) + "=" + url.QueryEscape(value) + "&"
	}
	return this.URLString+paramStr
}

func (this *Client) createGetRequest() (*http.Request, error) {
	return http.NewRequest(this.Method, this.createGetURL(), nil)
}

func (this *Client) createPostRequest() (*http.Request, error) {
	var param = make(url.Values)
	for key, value := range this.params {
		param[key] = []string{value}
	}

	if _, ok := this.headers["Content-Type"]; !ok {
		this.headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
	return http.NewRequest(this.Method, this.URLString, strings.NewReader(param.Encode()))
}

func (this *Client) doRequest() (*http.Response, error) {
	if this.headers == nil {
		this.headers = make(map[string]string)
	}

	var method = strings.ToUpper(this.Method)

	var requst *http.Request
	var err error

	if method == K_HTTP_REQ_METHOD_GET || method == K_HTTP_REQ_METHOD_HEAD{
		requst, err = this.createGetRequest()
	} else {
		requst, err = this.createPostRequest()
	}

	if err != nil {
		return nil, err
	}

	for key, value := range this.headers {
		requst.Header.Set(key, value)
	}

	var client = &http.Client{}
	return client.Do(requst)
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


func DoGet(reqUrl string, params map[string]string) (map[string]interface{}, error) {
	var paramStr = "?"
	for key, value := range params {
		paramStr = paramStr + url.QueryEscape(key) + "=" + url.QueryEscape(value) + "&"
	}
	responseHandler, err := http.Get(reqUrl + paramStr)
	if err != nil {
		return nil, err
	}
	defer responseHandler.Body.Close()

	bodyByte, err := ioutil.ReadAll(responseHandler.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	err = json.Unmarshal(bodyByte, &result)
	return result, err
}

func DoPost(reqUrl string, params map[string]string) (map[string]interface{}, error) {
	var param = make(url.Values)
	for key, value := range params {
		param[key] = []string{value}
	}

	reqBody := bytes.NewBufferString(param.Encode())

	responseHandler, err := http.Post(reqUrl, "application/x-www-form-urlencoded", reqBody)
	if err != nil {
		return nil, err
	}
	defer responseHandler.Body.Close()

	bodyByte, err := ioutil.ReadAll(responseHandler.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	err = json.Unmarshal(bodyByte, &result)
	return result, err
}
