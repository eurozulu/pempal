package byteparsers

import (
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
)

var derPublicKeyFileExtensions = map[string]bool{"": true, "key": true, "pub": true}

type derPublicKeyParser struct {
}

func (d derPublicKeyParser) MatchPath(path string) bool {
	return derPublicKeyFileExtensions[strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))]
}

func (d derPublicKeyParser) Parse(data []byte) ([]*pem.Block, []byte, error) {
	_, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		_, err = x509.ParsePKCS1PublicKey(data)
		if err != nil {
			return nil, data, err
		}
	}
	return []*pem.Block{{
		Type:  "PUBLIC KEY",
		Bytes: data,
	}}, nil, nil
}
