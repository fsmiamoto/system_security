package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func GenKey() (string, error) {
	var key [4]byte
	_, err := rand.Read(key[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key[:]), nil
}

func Encrypt(key, iv, plainText []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padded := PKCS5Padding(plainText, block.BlockSize())
	ciphered := make([]byte, len(padded))

	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(ciphered, padded)

	return ciphered, nil
}

func Decrypt(key, iv, cipherText []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bm := cipher.NewCBCDecrypter(block, iv)
	unciphered := make([]byte, len(cipherText))

	bm.CryptBlocks(unciphered, cipherText)

	return PKCS5Unpadding(unciphered), nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5Unpadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
