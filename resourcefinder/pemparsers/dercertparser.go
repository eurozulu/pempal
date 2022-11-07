package pemparsers

import (
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
)

var derCertFileExtensions = []string{"", "crt", "cert", "der"}

type derCertificateParser struct {
}

func (d derCertificateParser) Match(path string) bool {
	return stringIndex(strings.ToLower(filepath.Ext(path)), derCertFileExtensions) >= 0
}

func (d derCertificateParser) Parse(data []byte) (PemBlocks, []byte, error) {
	certs, err := x509.ParseCertificates(data)
	if err != nil {
		return nil, data, err
	}
	var pems PemBlocks
	for _, cert := range certs {
		pems = append(pems, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
	}
	return pems, nil, nil
}
