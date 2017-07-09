package request

import (
	"net/http"
	"encoding/json"
)

type Response struct {
	*http.Response
	Data []byte
	Error error
}

func (this *Response) Bytes() ([]byte, error) {
	return this.Data, this.Error
}

func (this *Response) MustBytes() ([]byte) {
	return this.Data
}

func (this *Response) String() (string, error) {
	return string(this.Data), this.Error
}

func (this *Response) MustString() (string) {
	return string(this.Data)
}

func (this *Response) UnmarshalJSON(v interface{}) (error) {
	return json.Unmarshal(this.Data, v)
}
