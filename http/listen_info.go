package http

import (
	"fmt"
	"unicode"
)

type ListenInfo struct {
	Enable   bool   `json:"Enable"`
	Port     int    `json:"Port"`
	Protocol string `json:"Protocol"`
}

func (l *ListenInfo) ListenAddress() string {
	return fmt.Sprintf("0.0.0.0:%v", l.Port)
}

func (l *ListenInfo) LowercaseProtocol() string {
	var result string
	for _, v := range l.Protocol {
		result += string(unicode.ToLower(v))
	}
	return result
}
