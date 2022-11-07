package pemparsers

import (
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
)

var derPublicKeyFileExtensions = []string{"", "key", "pub"}

type derPublicKeyParser struct {
}

func (d derPublicKeyParser) Match(path string) bool {
	return stringIndex(strings.ToLower(filepath.Ext(path)), derPublicKeyFileExtensions) >= 0
}

func (d derPublicKeyParser) Parse(data []byte) (PemBlocks, []byte, error) {
	_, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, data, err
	}
	return PemBlocks{&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: data,
	}}, nil, nil
}
