package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
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

func Encrypt(key, iv, plainText []byte) (ciphered []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			ciphered = nil
			err = errors.New("encrypt: invalid key")
			return
		}
	}()

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padded := pKCS5Padding(plainText, block.BlockSize())
	ciphered = make([]byte, len(padded))

	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(ciphered, padded)

	return ciphered, nil
}

func Decrypt(key, iv, cipherText []byte) (plain []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			plain = nil
			err = errors.New("decrypt: invalid key")
			return
		}
	}()

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bm := cipher.NewCBCDecrypter(block, iv)
	plain = make([]byte, len(cipherText))

	bm.CryptBlocks(plain, cipherText)

	return pKCS5Unpadding(plain), nil
}

func pKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pKCS5Unpadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
