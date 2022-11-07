package pemparsers

import (
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
)

var derCSRFileExtensions = []string{"", "csr", "der"}

type derCSRParser struct {
}

func (d derCSRParser) Match(path string) bool {
	return stringIndex(strings.ToLower(filepath.Ext(path)), derCSRFileExtensions) >= 0
}

func (d derCSRParser) Parse(data []byte) (PemBlocks, []byte, error) {
	csr, err := x509.ParseCertificateRequest(data)
	if err != nil {
		return nil, data, err
	}
	return PemBlocks{&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr.Raw,
	}}, nil, nil
}
