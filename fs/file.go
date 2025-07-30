package fs

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path"
)

func CheckFileExist(directory string) bool {
	f, err := os.Stat(directory)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

// WriteFullpathFile comment
func WriteFullpathFile(fullpathFilename, content string) error {
	directory := path.Dir(fullpathFilename)
	filename := path.Base(fullpathFilename)
	return WriteFile(directory, filename, content)
}

// WriteFile comment
func WriteFile(directory, filename, content string) error {
	err := CheckDirectoryExistAndCreateIfNot(directory)
	if err != nil {
		return errors.New("[brahms-utils, WriteFile]check directory failed: " + err.Error())
	}
	fullpathFilename := path.Join(directory, filename)
	f, err := os.OpenFile(fullpathFilename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("[brahms-utils, WriteFile]open file failed: " + err.Error())
	}
	f.WriteString(content)
	f.Close()
	return nil
}

// GetFileReader comment
func GetFileReader(fullpathFilename string) (*os.File, *bufio.Reader, error) {
	_, err := os.Stat(fullpathFilename)
	if err != nil {
		return nil, nil, err
	}
	f, err := os.Open(fullpathFilename)
	if err != nil {
		return nil, nil, err
	}
	return f, bufio.NewReader(f), nil
}

// CompareFile comment
func CompareFile(fullpathFilename1, fullpathFilename2 string) (bool, error) {
	hash1, err := GetFileHashSHA256(fullpathFilename1)
	if err != nil {
		return false, err
	}
	hash2, err := GetFileHashSHA256(fullpathFilename2)
	if err != nil {
		return false, err
	}
	if hash1 == hash2 {
		return true, nil
	}
	return false, nil
}

// CompareFileAndContent comment
func CompareFileAndContent(fullpathFilename, content string) (bool, error) {
	contentBytes, err := os.ReadFile(fullpathFilename)
	if err != nil {
		return false, err
	}
	if string(contentBytes) == content {
		return true, nil
	}
	return false, nil
}

// CopyFile comment
func CopyFile(sourceFullpathFilename, destFullpathFilename string) (written int64, err error) {
	src, err := os.Open(sourceFullpathFilename)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(destFullpathFilename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func GetFileSize(filename string) (int64, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	var size int64
	buf := make([]byte, 10240)
	for {
		length, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				size += int64(length)
				break
			} else {
				return 0, err
			}
		} else {
			size += int64(length)
		}
	}
	return size, nil
}

func GetFileHashSHA256(filename string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	var content []byte
	buf := make([]byte, 10240)
	for {
		length, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				content = append(content, buf[:length]...)
				break
			} else {
				return "", err
			}
		} else {
			content = append(content, buf[:length]...)
		}
	}
	hash := sha256.Sum256(content)
	return base64.URLEncoding.EncodeToString(hash[:]), nil
}

func GetContentHashSHA256(contentBytes []byte) string {
	hash := sha256.Sum256(contentBytes)
	return base64.URLEncoding.EncodeToString(hash[:])
}

func GetFileSizeAndHashSHA256(filename string) (int64, string, error) {
	size, err := GetFileSize(filename)
	if err != nil {
		return 0, "", err
	}
	hash, err := GetFileHashSHA256(filename)
	if err != nil {
		return 0, "", err
	}
	return size, hash, nil
}
