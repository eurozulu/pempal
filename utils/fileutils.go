package utils

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

func FileOrDirectoryExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func FileExists(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

func DirectoryExists(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func EnsurePathExists(path string, perm os.FileMode) error {
	if DirectoryExists(path) {
		return nil
	}
	return os.MkdirAll(path, perm)
}

func ReadStdIn() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		buf.Write(scan.Bytes())
	}
	if scan.Err() != nil && scan.Err() != io.EOF {
		return nil, scan.Err()
	}
	return buf.Bytes(), nil
}
