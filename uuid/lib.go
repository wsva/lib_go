package uuid

import (
	"fmt"
	mathrand "math/rand"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func New() string {
	return fmt.Sprint(uuid.NewV4())
}

func NewDate() string {
	t := time.Now()
	timestr := t.Format("20060102150405")
	generator := mathrand.New(mathrand.NewSource(t.UnixNano()))
	randstr := fmt.Sprintf("%v%v%v",
		strings.Repeat("0", 32), generator.Int63n(1000000), generator.Int63n(1000000))
	return fmt.Sprintf("%v%v",
		timestr, randstr[len(randstr)-(32-len(timestr)):])
}
