package pemparsers

import (
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
)

var derCRLFileExtensions = []string{"", "crl", "der"}

type derCRLParser struct {
}

func (d derCRLParser) Match(path string) bool {
	return stringIndex(strings.ToLower(filepath.Ext(path)), derCRLFileExtensions) >= 0
}

func (d derCRLParser) Parse(data []byte) (PemBlocks, []byte, error) {
	crl, err := x509.ParseRevocationList(data)
	if err != nil {
		return nil, data, err
	}
	return PemBlocks{&pem.Block{
		Type:  "X509 CRL",
		Bytes: crl.Raw,
	}}, nil, nil
}
