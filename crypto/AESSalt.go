package crypto

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
)

const (
	AES256SaltIdentifier = "{AES256Sa%vlt}"
	AES128SaltIdentifier = "{AES128Sa%vlt}"
)

func AES256SaltEncrypt(key, iv string, text string) (string, error) {
	salt := fmt.Sprint(RandomInt64(10000) + 1)
	ctext, err := aesEncrypt(key+salt, iv+salt, AES256KeySize, []byte(text))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(AES256SaltIdentifier, salt) +
		base64.URLEncoding.EncodeToString(ctext), nil
}

func AES256SaltDecrypt(key, iv string, ctext string) (string, error) {
	ctextBase64, salt, ok := ParseAES256SaltText(ctext)
	if !ok {
		return "", errors.New("not aes cipher text")
	}
	ctextBytes, err := base64.URLEncoding.DecodeString(ctextBase64)
	if err != nil {
		return "", err
	}
	text, err := aesDecrpt(key+salt, iv+salt, AES256KeySize, ctextBytes)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func AES128SaltEncrypt(key, iv string, text string) (string, error) {
	salt := fmt.Sprint(RandomInt64(10000) + 1)
	ctext, err := aesEncrypt(key+salt, iv+salt, AES128KeySize, []byte(text))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(AES128SaltIdentifier, salt) +
		base64.URLEncoding.EncodeToString(ctext), nil
}

func AES128SaltDecrypt(key, iv string, ctext string) (string, error) {
	ctextBase64, salt, ok := ParseAES128SaltText(ctext)
	if !ok {
		return "", errors.New("not aes cipher text")
	}
	ctextBytes, err := base64.URLEncoding.DecodeString(ctextBase64)
	if err != nil {
		return "", err
	}
	text, err := aesDecrpt(key+salt, iv+salt, AES128KeySize, ctextBytes)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func ParseAES256SaltText(cetxt string) (string, string, bool) {
	reg := regexp.MustCompile("^" + fmt.Sprintf(AES256SaltIdentifier, `(\d+)`))
	submatch := reg.FindStringSubmatch(cetxt)
	if len(submatch) != 2 {
		return "", "", false
	}
	return reg.ReplaceAllString(cetxt, ""), submatch[1], true
}

func ParseAES128SaltText(cetxt string) (string, string, bool) {
	reg := regexp.MustCompile("^" + fmt.Sprintf(AES128SaltIdentifier, `(\d+)`))
	submatch := reg.FindStringSubmatch(cetxt)
	if len(submatch) != 2 {
		return "", "", false
	}
	return reg.ReplaceAllString(cetxt, ""), submatch[1], true
}

func IsAES256SaltText(cetxt string) bool {
	reg := regexp.MustCompile("^" + fmt.Sprintf(AES256SaltIdentifier, `(\d+)`))
	return reg.MatchString(cetxt)
}

func IsAES128SaltText(cetxt string) bool {
	reg := regexp.MustCompile("^" + fmt.Sprintf(AES256SaltIdentifier, `(\d+)`))
	return reg.MatchString(cetxt)
}
