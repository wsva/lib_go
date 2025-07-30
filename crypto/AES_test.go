package crypto_test

import (
	"fmt"
	"testing"

	"github.com/wsva/lib_go/crypto"
)

func TestAES128(T *testing.T) {
	ctext, err := crypto.AES128Encrypt("1", "2", "3")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ctext)
}
