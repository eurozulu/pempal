package pemio

import (
	"fmt"
	"strings"

	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/pkcs12"
)

func AllPemParsers() []PemParser {
	return []PemParser{
		&pemParser{},
		&pkcs12Parser{},
		&derParser{},
	}
}

func NewPemParser(readertype string) PemParser {
	switch readertype {
	case "pem":
		return &pemParser{}
	case "der":
		return &derParser{}
	case "p12":
		return &pkcs12Parser{}
	default:
		return nil
	}
}

type PemParser interface {
	ParsePems(by []byte) ([]*pem.Block, error)
}

type pemParser struct {
}

func (pr pemParser) ParsePems(by []byte) ([]*pem.Block, error) {
	s := string(by)
	if !strings.Contains(s, "-----BEGIN") || !strings.Contains(s, "-----END") {
		return nil, fmt.Errorf("unknown format")
	}
	var blks []*pem.Block
	data := by
	for {
		b, r := pem.Decode(data)
		if b == nil {
			if len(r) == len(by) {
				return nil, fmt.Errorf("no pems found")
			}
			break
		}
		blks = append(blks, b)
		data = r
	}
	return blks, nil
}

type derParser struct {
}

func (pr derParser) ParsePems(by []byte) ([]*pem.Block, error) {
	cs, err := x509.ParseCertificates(by)
	if err == nil {
		var blks []*pem.Block
		for _, c := range cs {
			blks = append(blks, &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: c.Raw,
			})
		}
		return blks, nil
	}
	_, err = x509.ParseCertificateRequest(by)
	if err == nil {
		return []*pem.Block{{
			Type:  "CERTIFICATE REQUEST",
			Bytes: by,
		}}, nil
	}

	_, err = x509.ParsePKCS8PrivateKey(by)
	if err == nil {
		return []*pem.Block{{
			Type:  "PRIVATE KEY",
			Bytes: by,
		}}, nil
	}

	_, err = x509.ParsePKIXPublicKey(by)
	if err == nil {
		return []*pem.Block{{
			Type:  "PUBLIC KEY",
			Bytes: by,
		}}, nil
	}

	return nil, fmt.Errorf("unknown format")
}

type pkcs12Parser struct {
	Password string
}

func (pr pkcs12Parser) ParsePems(by []byte) ([]*pem.Block, error) {
	return pkcs12.ToPEM(by, pr.Password)
}
