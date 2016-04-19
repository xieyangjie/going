package encoding

import (
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
)

func packageData(originalData []byte, packageSize int) (r [][]byte) {
	var length = len(originalData)
	var diff = length % packageSize
	var count = length / packageSize;
	if diff > 0 {
		count += 1
	}

	r = make([][]byte, 0, count)
	var endIndex = packageSize;
	for i:=0; i<count; i++ {
		if (i == count - 1 && diff > 0) {
			endIndex = i * packageSize + diff
		} else {
			endIndex = (i + 1) * packageSize
		}
		var b = originalData[i * packageSize : endIndex]
		r = append(r, b)
	}
	return r
}

func RSAEncrypt(plaintext, key []byte) ([]byte, error) {
	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return nil, nil
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	var pub = pubInterface.(*rsa.PublicKey)

	var datas = packageData(plaintext, 128 - 11)
	var cipherDatas []byte = make([]byte, 0, 0)

	for _, d := range datas {
		var c, e = rsa.EncryptPKCS1v15(rand.Reader, pub, d)
		if e != nil {
			return nil, e
		}
		cipherDatas = append(cipherDatas, c...)
	}

	return cipherDatas, nil
}

func RSADecrypt(ciphertext, key []byte) ([]byte, error) {
	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return nil, nil
	}

	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var datas = packageData(ciphertext, 128)
	var plainDatas []byte = make([]byte, 0, 0)

	for _, d := range datas {
		var p, e = rsa.DecryptPKCS1v15(rand.Reader, pri, d)
		if e != nil {
			return nil, e
		}
		plainDatas = append(plainDatas, p...)
	}
	return plainDatas, nil
}