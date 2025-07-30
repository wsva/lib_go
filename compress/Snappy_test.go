package compress_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/wsva/lib_go/compress"
)

func TestSnappy(T *testing.T) {
	fmt.Println("TestSnappy")
	text := strings.Repeat("12345678901234567890", 300)
	ctext := compress.SnappyCompress(text)
	fmt.Println(ctext)
	dtext, _ := compress.SnappyDecompress(ctext)
	fmt.Println(dtext)
}
