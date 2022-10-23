package resources

import (
	"crypto/x509"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var resourceFileTypes = []string{"pem", "der", "crt", "cer", "cert", "csr"}

type ResourceParser interface {
	AcceptFile(path string, d fs.DirEntry) bool
	Parse(p string) (Resource, error)
}

type certificateParser struct {
}

func (c certificateParser) AcceptFile(path string, d fs.DirEntry) bool {
	e := strings.ToLower(filepath.Ext(d.Name()))
	if !d.IsDir() && containsString(e, resourceFileTypes) {
		return true
	}
	return false
}

func (c certificateParser) Parse(p string) (Resource, error) {
	by, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	crt, err := x509.ParseCertificate(by)
	if err != nil {
		return nil, err
	}

}

func containsString(s string, ss []string) bool {
	for _, sz := range ss {
		if s == sz {
			return true
		}
	}
	return false
}
