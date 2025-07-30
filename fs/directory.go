package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

/*
如果是软连接，则继续追踪是否是目录
*/
func IsDir(directory string, entry fs.DirEntry) bool {
	if entry.IsDir() {
		return true
	}
	if entry.Type()&fs.ModeSymlink == 0 {
		return false
	}
	fullpath := filepath.Join(directory, entry.Name())
	realpath, err := filepath.EvalSymlinks(fullpath)
	if err != nil {
		//invalid link
		return false
	}
	fi, err := os.Stat(realpath)
	if err != nil {
		//invalid link
		return false
	}
	return fi.IsDir()
}

func CheckDirectoryExist(directory string) bool {
	f, err := os.Stat(directory)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func CheckDirectoryExistAndCreateIfNot(directory string) error {
	if CheckDirectoryExist(directory) {
		return nil
	}
	if CheckFileExist(directory) {
		return errors.New(directory + " is a file")
	}
	if !CheckDirectoryExist(directory) {
		os.MkdirAll(directory, 0777)
	}
	if !CheckDirectoryExist(directory) {
		return fmt.Errorf("create direcory %s error", directory)
	}
	return nil
}

func CheckPathAUnderB(A, B string) bool {
	rel, err := filepath.Rel(B, A)
	if err != nil {
		return false
	}
	return !strings.Contains(rel, "..")
}

// GetExecutableFullpath comment
func GetExecutableFullpath() (string, error) {
	ePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ePath), nil
}

func Basepath() string {
	ePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ePath)
}

// GetExecutableName comment
func GetExecutableName() (string, error) {
	ePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Base(ePath), nil
}
