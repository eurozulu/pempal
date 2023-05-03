package templates

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultPermission = 0640

type BlobStore interface {
	Contains(name string) bool
	Read(name string) ([]byte, error)
	Write(name string, blob []byte) error
	Delete(name string) error
}

type fileBlobStore struct {
	root       string
	extensions map[string]bool
}

func (b fileBlobStore) Contains(name string) bool {
	return len(b.fileNames(name)) > 0
}

func (b fileBlobStore) Read(name string) ([]byte, error) {
	p, err := b.resolveNameMustExist(name)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b fileBlobStore) Write(name string, blob []byte) error {
	p := b.resolveName(name)
	if fileExists(p, false) {
		return os.ErrExist
	}
	return os.WriteFile(p, blob, defaultPermission)
}

func (b fileBlobStore) Delete(name string) error {
	return os.Remove(b.resolveName(name))
}

func (b fileBlobStore) fileNames(findFirst ...string) []string {
	des, err := os.ReadDir(b.root)
	if err != nil {
		log.Printf("failed to read template root path '%s',  %v", b.root, err)
		return nil
	}
	var names []string
	var findingFirst string
	if len(findFirst) > 0 {
		findingFirst = strings.ToLower(findFirst[0])
	}
	for _, de := range des {
		if strings.HasPrefix(de.Name(), ".") || de.IsDir() {
			continue
		}
		if len(b.extensions) > 0 {
			ext := strings.ToLower(strings.TrimLeft(filepath.Ext(de.Name()), "."))
			if !b.extensions[ext] {
				continue
			}
		}
		if findingFirst != "" && strings.HasPrefix(strings.ToLower(de.Name()), findingFirst) {
			return []string{de.Name()}
		}
		names = append(names, de.Name())
	}
	if findingFirst != "" {
		// Not found what we're looking for
		return nil
	}
	return names
}

func (b fileBlobStore) resolveNameMustExist(name string) (string, error) {
	p := b.resolveName(name)
	if !fileExists(p, false) {
		return "", os.ErrNotExist
	}
	return p, nil
}

func (b fileBlobStore) resolveName(name string) string {
	return filepath.Join(b.root, name)
}

func fileExists(path string, asDirectory bool) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir() == asDirectory
}

func NewFileBlobStore(rootpath string, extensions ...string) (BlobStore, error) {
	if !fileExists(rootpath, true) {
		return nil, fmt.Errorf("root path '%s' does not exist", rootpath)
	}
	var m map[string]bool
	if len(extensions) > 0 {
		m = map[string]bool{}
		for _, ex := range extensions {
			m[strings.TrimSpace(strings.ToLower(ex))] = true
		}
	}
	return &fileBlobStore{
		root:       rootpath,
		extensions: m,
	}, nil
}
