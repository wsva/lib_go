package ant_design_pro

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/wsva/lib_go/net"
)

type HttpReq struct {
	IP       string          `json:"ip,omitempty"`
	Username string          `json:"username,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
}

func (s *HttpReq) Reader() *bytes.Reader {
	jsonBytes, _ := json.Marshal(*s)
	return bytes.NewReader(jsonBytes)
}

func ParseHttpReq(r *http.Request, limit int64) (*HttpReq, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(io.LimitReader(r.Body, limit))
	if err != nil {
		return nil, err
	}
	var req HttpReq
	json.Unmarshal(body, &req)
	req.IP = net.GetIPFromRequest(r).String()
	return &req, nil
}

func NewHttpReq(ip, user string, data interface{}) *HttpReq {
	jsonBytes, _ := json.Marshal(data)
	return &HttpReq{
		IP:       ip,
		Username: user,
		Data:     json.RawMessage(jsonBytes),
	}
}
