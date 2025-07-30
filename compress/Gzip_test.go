package compress_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wsva/lib_go/compress"
)

func TestGzip(T *testing.T) {
	fmt.Println("TestGzip")
	text := strings.Repeat("12345678901234567890", 300)
	ctext, err := compress.GzipCompress(text)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ctext)
	dtext, err := compress.GzipDecompress(ctext)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dtext)
}
