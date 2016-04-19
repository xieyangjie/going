package encoding

import (
	"bytes"
)

func ZeroPadding(text []byte, blockSize int) []byte {
	var diff = blockSize - len(text) % blockSize
	var paddingText = bytes.Repeat([]byte{0}, diff)
	return append(text, paddingText...)
}


func PKCS7Padding(text []byte, blockSize int) []byte {
	var diff = blockSize - len(text) % blockSize
	var paddingText = bytes.Repeat([]byte{byte(diff)}, diff)
	return append(text, paddingText...)
}
