package encoding

import (
	"encoding/base64"
)

func Base64Encoding(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func Base64Decoding(value string) string {
	if bytes, err := base64.StdEncoding.DecodeString(value); err == nil {
		return string(bytes)
	}
	return ""
}