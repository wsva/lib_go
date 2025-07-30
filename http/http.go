package http

import (
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type HttpClient struct {
	//address of https server
	Address string

	//http.MethodPost, http.MethodGet
	Method string

	//used in POST method
	Data io.Reader

	Timeout time.Duration // second

	CookieList []*http.Cookie

	HeaderMap map[string]string

	//used to limit the response size to read
	LimitResponse bool //default false
	LimitBytes    int64
}

func (h *HttpClient) getHttpClient() (*http.Client, error) {
	timeout := h.Timeout
	if timeout == 0 {
		timeout = 10
	}
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	return &http.Client{
		Transport: tr,
		Timeout:   timeout * time.Second,
	}, nil
}

func (h *HttpClient) newRequest() (*http.Request, error) {
	var request *http.Request
	var err error
	switch h.Method {
	case http.MethodGet:
		request, err = http.NewRequest(h.Method, h.Address, nil)
		if err != nil {
			return nil, err
		}
	case http.MethodPost:
		request, err = http.NewRequest(http.MethodPost, h.Address, h.Data)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported method")
	}
	for _, v := range h.CookieList {
		request.AddCookie(v)
	}
	for k, v := range h.HeaderMap {
		request.Header.Set(k, v)
	}
	return request, nil
}

func (h *HttpClient) DoRequest() ([]byte, error) {
	client, err := h.getHttpClient()
	if err != nil {
		return nil, err
	}
	request, err := h.newRequest()
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if h.LimitResponse {
		return io.ReadAll(io.LimitReader(resp.Body, h.LimitBytes))
	} else {
		return io.ReadAll(resp.Body)
	}
}

func (h *HttpClient) DoRequestRaw() (*http.Response, error) {
	client, err := h.getHttpClient()
	if err != nil {
		return nil, err
	}
	request, err := h.newRequest()
	if err != nil {
		return nil, err
	}
	return client.Do(request)
}
