package byteparsers

import (
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
)

var derCRLFileExtensions = map[string]bool{"": true, "crl": true, "der": true}

type derCRLParser struct {
}

func (d derCRLParser) MatchPath(path string) bool {
	return derCRLFileExtensions[strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))]
}

func (d derCRLParser) Parse(data []byte) ([]*pem.Block, []byte, error) {
	crl, err := x509.ParseRevocationList(data)
	if err != nil {
		return nil, data, err
	}
	return []*pem.Block{{
		Type:  "X509 CRL",
		Bytes: crl.Raw,
	}}, nil, nil
}
