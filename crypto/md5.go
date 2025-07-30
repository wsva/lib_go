package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5HexString(src string) string {
	digest := md5.Sum([]byte(src))
	return hex.EncodeToString(digest[:])
}
