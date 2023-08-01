package templates

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultPermission = 0640

var StoreFilePermission = os.FileMode(defaultPermission)

type fileStore struct {
	root          string
	fileExtension string
}

func (b fileStore) Contains(name string) bool {
	return len(b.fileNames(name)) > 0
}

func (b fileStore) Read(name string) ([]byte, error) {
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

func (b fileStore) Write(name string, blob []byte) error {
	p := b.resolveName(name)
	if fileExists(p, false) {
		return os.ErrExist
	}
	return os.WriteFile(p, blob, StoreFilePermission)
}

func (b fileStore) Delete(name string) error {
	return os.Remove(b.resolveName(name))
}

func (b fileStore) Names() []string {
	return b.fileNames()
}

func (b fileStore) fileNames(names ...string) []string {
	des, err := os.ReadDir(b.root)
	if err != nil {
		log.Printf("failed to read template root path '%s',  %v", b.root, err)
		return nil
	}
	names = sliceToLower(names)
	var found []string

	for _, de := range des {
		// ignore subdirs && hidden files
		if strings.HasPrefix(de.Name(), ".") || de.IsDir() {
			continue
		}
		if b.fileExtension != "" && strings.ToLower(strings.TrimLeft(filepath.Ext(de.Name()), ".")) != b.fileExtension {
			continue
		}
		// if looking for names, is it a name we're interested in
		if len(names) > 0 && !containsString(stripeFileExtension(de.Name()), names) {
			continue
		}
		found = append(found, de.Name())
		if len(names) == 1 && len(found) == 1 {
			break
		}
	}
	return found
}

func (b fileStore) resolveNameMustExist(name string) (string, error) {
	p := b.resolveName(name)
	if !fileExists(p, false) {
		return "", os.ErrNotExist
	}
	return p, nil
}

func (b fileStore) resolveName(name string) string {
	if filepath.Ext(name) == "" && b.fileExtension != "" {
		name = strings.Join([]string{name, b.fileExtension}, ".")
	}
	return filepath.Join(b.root, name)
}

func fileExists(path string, asDirectory bool) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir() == asDirectory
}

func containsString(s string, ss []string) bool {
	for _, sz := range ss {
		if sz == s {
			return true
		}
	}
	return false
}

func stripeFileExtension(s string) string {
	if ext := filepath.Ext(s); ext != "" {
		return s[:len(s)-len(ext)]
	}
	return s
}

func sliceToLower(ss []string) []string {
	found := make([]string, len(ss))
	for i, s := range ss {
		found[i] = strings.ToLower(s)
	}
	return found
}

func newFileStore(rootpath string, fileExtension string) (ByteStore, error) {
	if !fileExists(rootpath, true) {
		return nil, fmt.Errorf("root path '%s' does not exist", rootpath)
	}
	return &fileStore{
		root:          rootpath,
		fileExtension: strings.TrimLeft(fileExtension, "."),
	}, nil
}
