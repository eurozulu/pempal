package byteparsers

import (
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"strings"
)

var derKeyFileExtensions = map[string]bool{"": true, "key": true, "prk": true}

type derPrivateKeyParser struct {
}

func (d derPrivateKeyParser) MatchPath(path string) bool {
	return derKeyFileExtensions[strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))]
}

func (d derPrivateKeyParser) Parse(data []byte) ([]*pem.Block, []byte, error) {
	_, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return nil, data, err
	}
	return []*pem.Block{{
		Type:  "PRIVATE KEY:",
		Bytes: data,
	}}, nil, nil
}
