package net

import (
	"fmt"
	"io"
	"net/http"
)

// copyHeader copies headers from one http.Header to another.
// http://golang.org/src/pkg/net/http/httputil/reverseproxy.go#L72
func CopyHeader(dst http.Header, src http.Header) {
	for k := range dst {
		dst.Del(k)
	}
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func GetRequestString(req *http.Request) string {
	path := req.URL.Path
	qurey := req.URL.Query().Encode()
	if qurey != "" {
		path += "?" + qurey
	}
	result := fmt.Sprintf("%v %v HTTP/1.1\r\n", req.Method, path)

	for k1, v1 := range req.Header {
		for _, v2 := range v1 {
			result += fmt.Sprintf("%v: %v\r\n", k1, v2)
		}
	}

	result += "\r\n"

	body, _ := io.ReadAll(req.Body)
	result += string(body)

	return result
}
