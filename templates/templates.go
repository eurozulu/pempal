package templates

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
)

type Template interface {
	String() string
}

func NewTemplate(pb *pem.Block) (Template, error) {
	switch pb.Type {
	case "CERTIFICATE":
		return NewCertificateTemplate(pb)

	case "CERTIFICATE REQUEST":
		return NewRequestTemplate(pb)

	case "PRIVATE KEY":
		return NewPrivateKeyTemplate(pb), nil

	case "PUBLIC KEY":
		k, err := x509.ParsePKIXPublicKey(pb.Bytes)
		if err != nil {
			return nil, err
		}
		return PublicKeyTemplate{key: k}, nil

	default:
		return nil, fmt.Errorf("%s is an unknown pem type", pb.Type)
	}
}

func NewTemplates(pbs []*pem.Block) ([]Template, error) {
	var tps []Template
	for _, pb := range pbs {
		t, err := NewTemplate(pb)
		if err != nil {
			return nil, err
		}
		tps = append(tps, t)
	}
	return tps, nil
}

func ApplyTemplate(dst Template, src Template) error {
	by, err := yaml.Marshal(src)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(by, dst)
}

func ApplyTemplates(dst Template, templates ...Template) error {
	for i := len(templates) - 1; i >= 0; i++ {
		for _, t := range templates {
			if err := ApplyTemplate(dst, t); err != nil {
				return err
			}
		}
	}
	return nil
}
