package io

import (
	"bytes"
	"strings"
	"sync"
)

type SingleWriter struct {
	b bytes.Buffer
	l sync.Mutex
}

func (sw *SingleWriter) Write(p []byte) (n int, err error) {
	sw.l.Lock()
	n, err = sw.b.Write(p)
	sw.l.Unlock()
	return
}

func (sw *SingleWriter) WriteString(s string) (n int, err error) {
	sw.l.Lock()
	n, err = sw.b.WriteString(s)
	sw.l.Unlock()
	return
}

func (sw *SingleWriter) Bytes() []byte {
	return sw.b.Bytes()
}

func (sw *SingleWriter) String() string {
	var builder strings.Builder
	builder.Write(sw.Bytes())
	return builder.String()
}
