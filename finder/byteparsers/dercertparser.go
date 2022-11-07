package byteparsers

import (
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
)

var derCertFileExtensions = map[string]bool{"": true, "crt": true, "cer": true, "cert": true, "der": true}

type derCertificateParser struct {
}

func (d derCertificateParser) MatchPath(path string) bool {
	return derCertFileExtensions[strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))]
}

func (d derCertificateParser) Parse(data []byte) ([]*pem.Block, []byte, error) {
	certs, err := x509.ParseCertificates(data)
	if err != nil {
		return nil, data, err
	}
	var pems []*pem.Block
	for _, cert := range certs {
		pems = append(pems, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
	}
	return pems, nil, nil
}
