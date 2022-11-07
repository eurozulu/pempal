package filepaths

import (
	"os"
	"path/filepath"
)

const defaultHome = "${HOME}/.pempal"
const defaultDirPermission = 750
const defaultFilePermission = 640

var envHomePath = os.ExpandEnv(os.Getenv("PP_HOME"))

func ListFiles(path string) ([]string, error) {
	fp := filepath.Join(HomePath(), path)
	dirs, err := os.ReadDir(fp)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, d := range dirs {
		dn := filepath.Join(path, d.Name())
		if !d.IsDir() {
			names = append(names, dn)
			continue
		}
		subNames, err := ListFiles(dn)
		if err != nil {
			return nil, err
		}
		for _, sn := range subNames {
			names = append(names, filepath.Join(dn, sn))
		}
	}
	return names, nil
}

func ReadFile(path string) ([]byte, error) {
	fp := filepath.Join(HomePath(), path)
	return os.ReadFile(fp)
}

func ContainsFile(path string) bool {
	return FileExists(filepath.Join(HomePath(), path))
}

func FileExists(fullpath string) bool {
	fi, err := os.Stat(fullpath)
	return err == nil && !fi.IsDir()
}
func DirExists(fullpath string) bool {
	fi, err := os.Stat(fullpath)
	return err == nil && fi.IsDir()
}

func WriteFile(path string, data []byte) error {
	fp := filepath.Join(HomePath(), path)
	if err := EnsureDirectory(filepath.Dir(fp)); err != nil {
		return err
	}
	return os.WriteFile(fp, data, defaultFilePermission)
}

func EnsureDirectory(path string) error {
	fp := filepath.Join(HomePath(), path)
	return os.MkdirAll(fp, defaultDirPermission)
}

func HomePath() string {
	if envHomePath == "" {
		return defaultHome
	}
	return envHomePath
}
