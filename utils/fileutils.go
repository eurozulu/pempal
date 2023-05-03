package utils

import "os"

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
