package templates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"gopkg.in/yaml.v2"
)

type Template interface {
	Name() string
}

type template struct {
	name string
	data []byte
}

func (t template) Name() string {
	return t.name
}

func (t template) MarshalYAML() (interface{}, error) {
	m := map[string]interface{}{}
	if err := yaml.Unmarshal(t.data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (t *template) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(t.data)
}

func TemplateAsMap(t Template) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	if err := ApplyTemplateToTarget(t, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func TemplateOfPem(blk *pem.Block) (Template, error) {
	switch model.ParseResourceTypeFromPEMType(blk.Type) {
	case model.PrivateKey:
		prk, err := x509.ParsePKCS8PrivateKey(blk.Bytes)
		if err != nil {
			return nil, err
		}
		return NewPrivateKeyTemplate(prk)
	case model.PublicKey:
		puk, err := x509.ParsePKIXPublicKey(blk.Bytes)
		if err != nil {
			return nil, err
		}
		return NewPublicKeyTemplate(puk)
	case model.Certificate:
		cert, err := x509.ParseCertificate(blk.Bytes)
		if err != nil {
			return nil, err
		}
		return NewCertificateTemplate(cert), nil
	case model.CertificateRequest:
		csr, err := x509.ParseCertificateRequest(blk.Bytes)
		if err != nil {
			return nil, err
		}
		return NewCSRTemplate(csr), nil
	case model.RevokationList:
		crl, err := x509.ParseRevocationList(blk.Bytes)
		if err != nil {
			return nil, err
		}
		return NewCRLTemplate(crl), nil
	default:
		return nil, fmt.Errorf("unknown pem type")
	}
}
