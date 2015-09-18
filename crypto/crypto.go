package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(value []byte) string {
	var m = md5.New()
	m.Write(value)
	var dst = hex.EncodeToString(m.Sum(nil))
	return dst
}