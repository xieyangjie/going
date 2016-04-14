package encoding

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func PKCS5Padding(text []byte, blockSize int) []byte {
	var diff = blockSize - len(text) % blockSize
	var paddingText = bytes.Repeat([]byte{byte(diff)}, diff)
	return append(text, paddingText...)
}

// AESCBCEncrypt 由key的长度决定是128, 192 还是 256
func AESCBCEncrypt(plaintext, key, iv []byte) ([]byte, error) {
	var block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	var blockSize = block.BlockSize()
	iv = iv[:blockSize]

	var text = make([]byte, len(plaintext))
	copy(text, plaintext)
	text = PKCS5Padding(text, blockSize)

	var mode = cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(text, text)
	return text, nil
}

func AESCBCDecrypt(ciphertext, key, iv []byte) ([]byte, error) {
	var block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	var blockSize = block.BlockSize()
	iv = iv[:blockSize]

	var text = make([]byte, len(ciphertext))

	var mode = cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(text, ciphertext)
	return text, nil
}


func AESCFBEncrypt(plaintext, key, iv []byte) ([]byte, error) {
	var block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var blockSize = block.BlockSize()
	iv = iv[:blockSize]

	var text = make([]byte, len(plaintext))

	var mode = cipher.NewCFBEncrypter(block, iv)
	mode.XORKeyStream(text, plaintext)
	return text ,nil
}

func AESCFBDecrypt(ciphertext, key, iv []byte) ([]byte, error) {
	var block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var blockSize = block.BlockSize()
	iv = iv[:blockSize]

	var text = make([]byte, len(ciphertext))

	var mode = cipher.NewCFBDecrypter(block, iv)
	mode.XORKeyStream(text, ciphertext)
	return text, nil
}