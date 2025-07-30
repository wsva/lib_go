package crypto

import (
	cryptorand "crypto/rand"
	"math/big"
	mathrand "math/rand"
	"time"
)

func RandomInt64(max int64) int64 {
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(max))
	if err != nil {
		r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
		return r.Int63n(max)
	}
	return n.Int64()
}
