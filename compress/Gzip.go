package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
)

func GzipCompressBytes(text []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(text)
	if err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GzipDecompressBytes(ctext []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(ctext))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(gz)
}

func GzipCompress(text string) (string, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte(text))
	if err != nil {
		return "", err
	}
	if err := gz.Flush(); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func GzipDecompress(ctext string) (string, error) {
	ctextBase64, err := base64.StdEncoding.DecodeString(ctext)
	if err != nil {
		return "", err
	}
	gz, err := gzip.NewReader(bytes.NewReader(ctextBase64))
	if err != nil {
		return "", err
	}
	text, err := io.ReadAll(gz)
	return string(text), err
}
