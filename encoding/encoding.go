package encoding

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(value []byte) []byte {
	var m = md5.New()
	m.Write(value)
	return m.Sum(nil)
}

func Md5String(value string) string {
	return hex.EncodeToString(Md5([]byte(value)))
}