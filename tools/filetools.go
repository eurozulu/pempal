package tools

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
func IsPathExistsFs(fsys fs.FS, path string) bool {
	_, err := fs.Stat(fsys, path)
	return err == nil
}

func IsFileExists(path string) bool {
	i, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !i.IsDir()
}

func IsFileExistsFS(fsys fs.FS, path string) bool {
	i, err := fs.Stat(fsys, path)
	if err != nil {
		return false
	}
	return !i.IsDir()
}

func IsDirExists(path string) bool {
	i, err := os.Stat(path)
	if err != nil {
		return false
	}
	return i.IsDir()
}
func IsDirExistsFS(fsys fs.FS, path string) bool {
	i, err := fs.Stat(fsys, path)
	if err != nil {
		return false
	}
	return i.IsDir()
}

func EnsureSecurePath(path string) error {
	if !IsPathExists(path) {
		if p := filepath.Dir(path); p != "." {
			if err := os.MkdirAll(path, 0777); err != nil {
				return err
			}
		}
		if err := os.MkdirAll(path, 0700); err != nil {
			return err
		}
	}
	s, _ := os.Stat(path)
	if s.Mode()&0077 != 0 {
		return fmt.Errorf("permissions (%s) on key directory %s are too open. No group or public access should be allowed", s.Mode().Perm(), path)
	}
	return nil
}
