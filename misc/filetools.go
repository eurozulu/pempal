package misc

import "os"

func FileExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}
func DirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
