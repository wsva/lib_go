package compress

import (
	"encoding/base64"

	"github.com/golang/snappy"
)

func SnappyCompress(text string) string {
	ctextBytes := snappy.Encode(nil, []byte(text))
	return base64.StdEncoding.EncodeToString(ctextBytes)
}

func SnappyDecompress(ctext string) (string, error) {
	text, err := base64.StdEncoding.DecodeString(ctext)
	if err != nil {
		return "", err
	}
	dest, err := snappy.Decode(nil, []byte(text))
	return string(dest), err
}
