package crypto

import (
	"crypto/md5"
)

func MD5(value []byte) []byte {
	var m = md5.New()
	m.Write(value)
	return m.Sum(nil)
}