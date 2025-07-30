package ftp

import (
	"io"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTP struct {
	IP       string
	Port     string
	Username string
	Password string
}

func (f *FTP) ReadFile(fullpathFilename string) ([]byte, error) {
	client, err := ftp.Dial(f.IP+":"+f.Port, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return nil, err
	}

	err = client.Login(f.Username, f.Password)
	if err != nil {
		return nil, err
	}

	r, err := client.Retr(fullpathFilename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

func (f *FTP) WriteFile(fullpathFilename string, reader io.Reader) error {
	client, err := ftp.Dial(f.IP+":"+f.Port, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return err
	}

	err = client.Login(f.Username, f.Password)
	if err != nil {
		return err
	}

	err = client.Stor(fullpathFilename, reader)
	if err != nil {
		return err
	}

	return nil
}
