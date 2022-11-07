package pemparsers

import (
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
)

var derKeyFileExtensions = []string{"", "key", "prk"}

type derPrivateKeyParser struct {
}

func (d derPrivateKeyParser) Match(path string) bool {
	return stringIndex(strings.ToLower(filepath.Ext(path)), derKeyFileExtensions) >= 0
}

func (d derPrivateKeyParser) Parse(data []byte) (PemBlocks, []byte, error) {
	_, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return nil, data, err
	}
	return PemBlocks{&pem.Block{
		Type:  "PRIVATE KEY:",
		Bytes: data,
	}}, nil, nil
}
