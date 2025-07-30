package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"regexp"
)

const (
	AES256KeySize = 256
	AES128KeySize = 128
)

const (
	AES256Identifier = "{AES256}"
	AES128Identifier = "{AES128}"
)

func getAESKeySHA256(key string) []byte {
	digest := sha256.Sum256([]byte(key))
	return digest[:]
}

func getAESKeyMD5(key string) []byte {
	digest := md5.Sum([]byte(key))
	return digest[:]
}

func getAESkey(key string, keysize int) ([]byte, error) {
	switch keysize {
	case AES256KeySize:
		return getAESKeySHA256(key), nil
	case AES128KeySize:
		return getAESKeyMD5(key), nil
	default:
		return nil, errors.New("wrong keysize")
	}
}

func aesEncrypt(key, iv string, keysize int, text []byte) ([]byte, error) {
	aeskey, err := getAESkey(key, keysize)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	aesiv, err := getAESkey(iv, blockSize*8)
	if err != nil {
		return nil, err
	}
	text = PKCS5Padding(text, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, aesiv)
	ctext := make([]byte, len(text))
	blockMode.CryptBlocks(ctext, text)
	return ctext, nil
}

func aesDecrpt(key, iv string, keysize int, ctext []byte) ([]byte, error) {
	aeskey, err := getAESkey(key, keysize)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	aesiv, err := getAESkey(iv, blockSize*8)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, aesiv)
	text := make([]byte, len(ctext))
	blockMode.CryptBlocks(text, ctext)
	text = PKCS5UnPadding(text)
	return text, nil
}

func PKCS5Padding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padtext...)
}

func PKCS5UnPadding(text []byte) []byte {
	length := len(text)
	unpadding := int(text[length-1])
	return text[:(length - unpadding)]
}

func AES256Encrypt(key, iv string, text string) (string, error) {
	ctext, err := aesEncrypt(key, iv, AES256KeySize, []byte(text))
	if err != nil {
		return "", err
	}
	return AES256Identifier + base64.URLEncoding.EncodeToString(ctext), nil
}

func AES256Decrypt(key, iv string, ctext string) (string, error) {
	ctextBase64, ok := ParseAES256Text(ctext)
	if !ok {
		return "", errors.New("not aes cipher text")
	}
	ctextBytes, err := base64.URLEncoding.DecodeString(ctextBase64)
	if err != nil {
		return "", err
	}
	text, err := aesDecrpt(key, iv, AES256KeySize, ctextBytes)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func AES128Encrypt(key, iv string, text string) (string, error) {
	ctext, err := aesEncrypt(key, iv, AES128KeySize, []byte(text))
	if err != nil {
		return "", err
	}
	return AES128Identifier + base64.URLEncoding.EncodeToString(ctext), nil
}

func AES128Decrypt(key, iv string, ctext string) (string, error) {
	ctextBase64, ok := ParseAES128Text(ctext)
	if !ok {
		return "", errors.New("not aes cipher text")
	}
	ctextBytes, err := base64.URLEncoding.DecodeString(ctextBase64)
	if err != nil {
		return "", err
	}
	text, err := aesDecrpt(key, iv, AES128KeySize, ctextBytes)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func ParseAES256Text(cetxt string) (string, bool) {
	reg := regexp.MustCompile("^" + AES256Identifier)
	if !reg.MatchString(cetxt) {
		return "", false
	}
	return reg.ReplaceAllString(cetxt, ""), true
}

func ParseAES128Text(cetxt string) (string, bool) {
	reg := regexp.MustCompile("^" + AES128Identifier)
	if !reg.MatchString(cetxt) {
		return "", false
	}
	return reg.ReplaceAllString(cetxt, ""), true
}

func IsAES256Text(cetxt string) bool {
	reg := regexp.MustCompile("^" + AES256Identifier)
	return reg.MatchString(cetxt)
}

func IsAES128Text(cetxt string) bool {
	reg := regexp.MustCompile("^" + AES128Identifier)
	return reg.MatchString(cetxt)
}
