package resourceio

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/model"
)

var ResourceParsers = []ResourceParser{
	PemResourceParser{},
	DerResourceParser{},
}

type ResourceParser interface {
	CanParse(data []byte) bool
	ParseResources(data []byte) ([]model.Resource, error)
}

type PemResourceParser struct {
}

func (p PemResourceParser) CanParse(data []byte) bool {
	i := bytes.Index(data, []byte("-----BEGIN "))
	if i < 0 && i+1 < len(data) {
		return false
	}
	return bytes.Index(data[i+1:], []byte("-----END ")) > 0
}

func (p PemResourceParser) ParseResources(data []byte) ([]model.Resource, error) {
	var parsed []model.Resource
	for len(data) > 0 {
		blk, rest := pem.Decode(data)
		if blk == nil {
			break
		}
		parsed = append(parsed, model.NewResource(blk))
		data = rest
	}
	return parsed, nil
}

type DerResourceParser struct{}

func (d DerResourceParser) CanParse(data []byte) bool {
	return findResourceType(data) != model.Unknown
}

func (d DerResourceParser) ParseResources(data []byte) ([]model.Resource, error) {
	t := findResourceType(data)
	if t == model.Unknown {
		return nil, fmt.Errorf("failed to parse as a FormatDER resource")
	}
	return []model.Resource{
		model.NewResource(
			&pem.Block{
				Type:  t.String(),
				Bytes: data,
			}),
	}, nil
}

func findResourceType(der []byte) model.ResourceType {
	_, err := x509.ParseCertificate(der)
	if err == nil {
		return model.Certificate
	}
	_, err = x509.ParseCertificateRequest(der)
	if err == nil {
		return model.CertificateRequest
	}
	_, err = x509.ParsePKIXPublicKey(der)
	if err == nil {
		return model.PublicKey
	}
	_, err = x509.ParsePKCS8PrivateKey(der)
	if err == nil {
		return model.PrivateKey
	}
	_, err = x509.ParseRevocationList(der)
	if err == nil {
		return model.RevokationList
	}
	return model.Unknown
}
