package crypto

import (
	"errors"
	"regexp"

	"github.com/wsva/lib_go/compress"
)

const (
	SA256Identifier = "{SA256}"
	SA128Identifier = "{SA128}"
	GA256Identifier = "{GA256}"
	GA128Identifier = "{GA128}"
)

// SA128Encode : Snappy + AES128
func SA128Encode(key, iv, text string) (string, error) {
	ctext := compress.SnappyCompress(text)
	ctext, err := AES128Encrypt(key, iv, ctext)
	if err != nil {
		return "", err
	}
	return SA128Identifier + ctext, nil
}

// SA128Decode : Snappy + AES128
func SA128Decode(key, iv, ctext string) (string, error) {
	text, ok := ParseSA128Text(ctext)
	if !ok {
		return "", errors.New("not SA128 text")
	}
	text, err := AES128Decrypt(key, iv, text)
	if err != nil {
		return "", err
	}
	return compress.SnappyDecompress(text)
}

// SA256Encode : Snappy + AES256
func SA256Encode(key, iv, text string) (string, error) {
	ctext := compress.SnappyCompress(text)
	ctext, err := AES256Encrypt(key, iv, ctext)
	if err != nil {
		return "", err
	}
	return SA256Identifier + ctext, nil
}

// SA256Decode : Snappy + AES256
func SA256Decode(key, iv, ctext string) (string, error) {
	text, ok := ParseSA256Text(ctext)
	if !ok {
		return "", errors.New("not SA256 text")
	}
	text, err := AES256Decrypt(key, iv, text)
	if err != nil {
		return "", err
	}
	return compress.SnappyDecompress(text)
}

// GA128Encode : Gzip + AES128
func GA128Encode(key, iv, text string) (string, error) {
	ctext, err := compress.GzipCompress(text)
	if err != nil {
		return "", err
	}
	ctext, err = AES128Encrypt(key, iv, ctext)
	if err != nil {
		return "", err
	}
	return GA128Identifier + ctext, err
}

// GA128Decode : Gzip + AES128
func GA128Decode(key, iv, ctext string) (string, error) {
	text, ok := ParseGA128Text(ctext)
	if !ok {
		return "", errors.New("not GA128 text")
	}
	text, err := AES128Decrypt(key, iv, text)
	if err != nil {
		return "", err
	}
	return compress.GzipDecompress(text)
}

// GA256Encode : Gzip + AES256
func GA256Encode(key, iv, text string) (string, error) {
	ctext, err := compress.GzipCompress(text)
	if err != nil {
		return "", err
	}
	ctext, err = AES256Encrypt(key, iv, ctext)
	if err != nil {
		return "", err
	}
	return GA256Identifier + ctext, err
}

// GA256Decode : Gzip + AES256
func GA256Decode(key, iv, ctext string) (string, error) {
	text, ok := ParseGA256Text(ctext)
	if !ok {
		return "", errors.New("not GA256 text")
	}
	text, err := AES256Decrypt(key, iv, text)
	if err != nil {
		return "", err
	}
	return compress.GzipDecompress(text)
}

func ParseSA256Text(cetxt string) (string, bool) {
	reg := regexp.MustCompile("^" + SA256Identifier)
	if !reg.MatchString(cetxt) {
		return "", false
	}
	return reg.ReplaceAllString(cetxt, ""), true
}

func ParseSA128Text(cetxt string) (string, bool) {
	reg := regexp.MustCompile("^" + SA128Identifier)
	if !reg.MatchString(cetxt) {
		return "", false
	}
	return reg.ReplaceAllString(cetxt, ""), true
}

func ParseGA256Text(cetxt string) (string, bool) {
	reg := regexp.MustCompile("^" + GA256Identifier)
	if !reg.MatchString(cetxt) {
		return "", false
	}
	return reg.ReplaceAllString(cetxt, ""), true
}

func ParseGA128Text(cetxt string) (string, bool) {
	reg := regexp.MustCompile("^" + GA128Identifier)
	if !reg.MatchString(cetxt) {
		return "", false
	}
	return reg.ReplaceAllString(cetxt, ""), true
}
