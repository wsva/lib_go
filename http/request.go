package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/wsva/lib_go/net"
)

type Request struct {
	IP       string          `json:"ip,omitempty"`
	Username string          `json:"username,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
}

func (r *Request) Reader() *bytes.Reader {
	jsonBytes, _ := json.Marshal(*r)
	return bytes.NewReader(jsonBytes)
}

func ParseRequest(r *http.Request, limit int64) (*Request, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(io.LimitReader(r.Body, limit))
	if err != nil {
		return nil, err
	}
	var req Request
	json.Unmarshal(body, &req)
	req.IP = net.GetIPFromRequest(r).String()
	return &req, nil
}

func NewRequest(ip, user string, data any) *Request {
	jsonBytes, _ := json.Marshal(data)
	return &Request{
		IP:       ip,
		Username: user,
		Data:     json.RawMessage(jsonBytes),
	}
}
