package pempal

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/pempal/templates"
)

func NewRequest(t *templates.RequestTemplate, k *templates.PrivateKeyTemplate) (*pem.Block, error) {
	if k == nil {
		return nil, fmt.Errorf("failed to read private key")
	}

	rt := &templates.RequestTemplate{}
	if err := templates.ApplyTemplate(rt, t); err != nil {
		return nil, err
	}
	r := x509.CertificateRequest(*rt)
	key, err := k.PrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to read private key  %v", err)
	}

	der, err := x509.CreateCertificateRequest(rand.Reader, &r, key)
	if err != nil {
		return nil, err
	}
	return EncodeRequest(der), nil
}

func EncodeRequest(der []byte) *pem.Block {
	return &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: der,
	}
}
