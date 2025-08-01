package location

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jlaffaye/ftp"

	wl_fs "github.com/wsva/lib_go/fs"
	wl_http "github.com/wsva/lib_go/http"
)

type LocationInterface interface {
	Download(dest string) error
	Upload(src string) error
}

type Location struct {
	Enable       bool            `json:"Enable"`
	LocationType string          `json:"LocationType"`
	LocationInfo json.RawMessage `json:"LocationInfo"`

	pointer LocationInterface `json:"-"`
}

func (l *Location) Parse() error {
	if !l.Enable {
		return errors.New("location is not ababled")
	}
	if l.pointer != nil {
		return nil
	}
	switch l.LocationType {
	case "Local", "Directory":
		var location LocationLocal
		err := json.Unmarshal(l.LocationInfo, &location)
		if err != nil {
			return err
		}
		l.pointer = &location
	case "WebHttp":
		var location LocationWebHttp
		err := json.Unmarshal(l.LocationInfo, &location)
		if err != nil {
			return err
		}
		l.pointer = &location
	case "WebHttps":
		var location LocationWebHttps
		err := json.Unmarshal(l.LocationInfo, &location)
		if err != nil {
			return err
		}
		l.pointer = &location
	case "FTP":
		var location LocationFTP
		err := json.Unmarshal(l.LocationInfo, &location)
		if err != nil {
			return err
		}
		l.pointer = &location
	default:
		return errors.New("unknown location type: " + l.LocationType)
	}
	return nil
}

func (l *Location) Download(dest string) error {
	err := l.Parse()
	if err != nil {
		return err
	}
	return l.pointer.Download(dest)
}

func (l *Location) Upload(src string) error {
	err := l.Parse()
	if err != nil {
		return err
	}
	return l.pointer.Upload(src)
}

type LocationLocal struct {
	Path string `json:"Path"`
}

// dest is fullpath filename of destination
func (l *LocationLocal) Download(dest string) error {
	_, sourceReader, err := wl_fs.GetFileReader(l.Path)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, sourceReader)
	if err != nil {
		return err
	}
	return nil
}

// src is fullpath filename of source
func (l *LocationLocal) Upload(src string) error {
	_, sourceReader, err := wl_fs.GetFileReader(src)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(l.Path)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, sourceReader)
	if err != nil {
		return err
	}
	return nil
}

type LocationFTP struct {
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Path     string `json:"Path"`
}

// dest is fullpath filename of destination
func (l *LocationFTP) Download(dest string) error {
	client, err := ftp.Dial(l.Host+":"+l.Port, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}
	err = client.Login(l.Username, l.Password)
	if err != nil {
		return err
	}
	r, err := client.Retr(l.Path)
	if err != nil {
		return err
	}
	defer r.Close()
	outputFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, r)
	if err != nil {
		return err
	}
	return nil
}

// src is fullpath filename of source
func (l *LocationFTP) Upload(src string) error {
	client, err := ftp.Dial(l.Host+":"+l.Port, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	err = client.Login(l.Username, l.Password)
	if err != nil {
		return err
	}

	file, reader, err := wl_fs.GetFileReader(src)
	if err != nil {
		return err
	}
	defer file.Close()

	err = client.Stor(l.Path, reader)
	if err != nil {
		return err
	}

	return nil
}

type LocationWebHttp struct {
	URL string `json:"URL"`
}

// dest is fullpath filename of destination
func (l *LocationWebHttp) Download(dest string) error {
	client := wl_http.HttpClient{
		Address: l.URL,
		Method:  http.MethodGet,
	}
	resp, err := client.DoRequest()
	if err != nil {
		return err
	}
	return wl_fs.WriteFullpathFile(dest, string(resp))
}

// src is fullpath filename of source
func (l *LocationWebHttp) Upload(src string) error {
	file, reader, err := wl_fs.GetFileReader(src)
	if err != nil {
		return err
	}
	defer file.Close()
	client := wl_http.HttpClient{
		Address: l.URL,
		Method:  http.MethodPost,
		Data:    reader,
	}
	_, err = client.DoRequest()
	return err
}

type LocationWebHttps struct {
	URL           string `json:"URL"`
	CACrtFile     string `json:"CACrtFile"`
	MutualTLS     bool   `json:"MutualTLS"`
	ClientCrtFile string `json:"ClientCrtFile"`
	ClientKeyFile string `json:"ClientKeyFile"`
}

// dest is fullpath filename of destination
func (l *LocationWebHttps) Download(dest string) error {
	client := wl_http.HttpsClient{
		ServerAddress: l.URL,
		Method:        http.MethodGet,
		CACrtFile:     l.CACrtFile,
		MutualTLS:     l.MutualTLS,
		ClientCrtFile: l.ClientCrtFile,
		ClientKeyFile: l.ClientKeyFile,
	}
	resp, err := client.DoRequest(false)
	if err != nil {
		return err
	}
	return wl_fs.WriteFullpathFile(dest, string(resp))
}

// src is fullpath filename of source
func (l *LocationWebHttps) Upload(src string) error {
	file, reader, err := wl_fs.GetFileReader(src)
	if err != nil {
		return err
	}
	defer file.Close()
	client := wl_http.HttpsClient{
		ServerAddress: l.URL,
		Method:        http.MethodPost,
		Data:          reader,
		CACrtFile:     l.CACrtFile,
		MutualTLS:     l.MutualTLS,
		ClientCrtFile: l.ClientCrtFile,
		ClientKeyFile: l.ClientKeyFile,
	}
	_, err = client.DoRequest(false)
	return err
}
