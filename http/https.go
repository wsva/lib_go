package http

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

type HttpsServer struct {
	Port string `json:"Port"`

	// used only by mutual mode
	// DO NOT set this field, when you don't want to use mutual https
	CACrtFile string `json:"CACrtFile"`

	ServerCrtFile string `json:"ServerCrtFile"`
	ServerKeyFile string `json:"ServerKeyFile"`
}

/*
ListenAndServeTLS is used for TLS authentication

在Oracle Linux 6系列服务器上，curl、wget都不支持TLS v1.0，只能支持SSLv3
但是因为安全问题，目前https服务普遍都不支持SSLv3了，导致在这些老操作系统上，不能使用https
nginx也是如此

如果一定要用curl访问，建议同时开启http服务
*/
func (s *HttpsServer) ListenAndServe(handler http.Handler) error {
	server := &http.Server{
		Addr:    ":" + s.Port,
		Handler: handler,
	}

	if len(s.CACrtFile) > 0 {
		caPool := x509.NewCertPool()
		crt, err := os.ReadFile(s.CACrtFile)
		if err != nil {
			return err
		}
		caPool.AppendCertsFromPEM(crt)
		/*
			当server要求client提供证书时，会将ClientCAs中的CA根证书发送给client
			client在本地所有证书中，寻找与CA根证书匹配的证书，如果能找到，浏览器就会弹出选择证书的窗口
			如果找不到，会报错：不接受您的登录证书,或者您可能没有提供登录证书
		*/
		server.TLSConfig = &tls.Config{
			ClientCAs:  caPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
	}

	server.SetKeepAlivesEnabled(false)
	return server.ListenAndServeTLS(s.ServerCrtFile, s.ServerKeyFile)
}

type HttpsClient struct {
	ServerAddress string

	//http.MethodPost, http.MethodGet
	Method string

	//used in POST method
	Data io.Reader

	//used to verify cetificates of https server
	CACrtFile string

	//used in case of mutual TLS authentication
	MutualTLS     bool
	ClientCrtFile string
	ClientKeyFile string

	Timeout time.Duration // second

	CookieList []*http.Cookie

	HeaderMap map[string]string

	//used to limit the response size to read
	LimitResponse bool //default false
	LimitBytes    int64
}

func (h *HttpsClient) getCACrtPool() (*x509.CertPool, error) {
	if h.CACrtFile == "" {
		return nil, errors.New("CACrt is empty")
	}
	caPool := x509.NewCertPool()
	contentBytes, err := os.ReadFile(h.CACrtFile)
	if err != nil {
		return nil, err
	}
	caPool.AppendCertsFromPEM(contentBytes)
	return caPool, nil
}

func (h *HttpsClient) getClientCert() (*tls.Certificate, error) {
	if h.ClientCrtFile == "" {
		return nil, errors.New("ClientCrt is empty")
	}
	if h.ClientKeyFile == "" {
		return nil, errors.New("ClientKey is empty")
	}
	clientCert, err := tls.LoadX509KeyPair(h.ClientCrtFile, h.ClientKeyFile)
	return &clientCert, err
}

func (h *HttpsClient) getHttpClient(skipVerify bool) (*http.Client, error) {
	caPool, err := h.getCACrtPool()
	if err != nil && !skipVerify {
		return nil, err
	}

	var tr *http.Transport
	if h.MutualTLS {
		clientCert, err := h.getClientCert()
		if err != nil {
			return nil, err
		}
		tr = &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				RootCAs:      caPool,
				Certificates: []tls.Certificate{*clientCert},
			},
		}
	} else {
		if skipVerify {
			tr = &http.Transport{
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
		} else {
			tr = &http.Transport{
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					RootCAs: caPool,
				},
			}
		}
	}

	timeout := h.Timeout
	if timeout == 0 {
		timeout = 10
	}
	return &http.Client{
		Transport: tr,
		Timeout:   timeout * time.Second,
	}, nil
}

func (h *HttpsClient) newRequest() (*http.Request, error) {
	var request *http.Request
	var err error
	switch h.Method {
	case http.MethodGet:
		request, err = http.NewRequest(h.Method, h.ServerAddress, nil)
		if err != nil {
			return nil, err
		}
	case http.MethodPost:
		request, err = http.NewRequest(http.MethodPost, h.ServerAddress, h.Data)
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

func (h *HttpsClient) DoRequest(skipVerify bool) ([]byte, error) {
	client, err := h.getHttpClient(skipVerify)
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

func (h *HttpsClient) DoRequestRaw(skipVerify bool) (*http.Response, error) {
	client, err := h.getHttpClient(skipVerify)
	if err != nil {
		return nil, err
	}
	request, err := h.newRequest()
	if err != nil {
		return nil, err
	}
	return client.Do(request)
}
