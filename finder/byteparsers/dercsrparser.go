package byteparsers

import (
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
)

var derCSRFileExtensions = map[string]bool{"": true, "csr": true, "der": true}

type derCSRParser struct {
}

func (d derCSRParser) MatchPath(path string) bool {
	return derCSRFileExtensions[strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))]
}

func (d derCSRParser) Parse(data []byte) ([]*pem.Block, []byte, error) {
	csr, err := x509.ParseCertificateRequest(data)
	if err != nil {
		return nil, data, err
	}
	return []*pem.Block{{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr.Raw,
	}}, nil, nil
}
