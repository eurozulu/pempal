package pempal

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/pempal/templates"
)

func NewCertificate(cert, ca *templates.CertificateTemplate, cakey *templates.PrivateKeyTemplate) (*pem.Block, error) {
	c := x509.Certificate(*cert)
	if c.PublicKey == nil {
		return nil, fmt.Errorf("new certificate has no public key")
	}
	cac := x509.Certificate(*ca)

	prk, err := cakey.PrivateKey()
	if err != nil {
		return nil, err
	}

	der, err := x509.CreateCertificate(rand.Reader, &c, &cac, c.PublicKey, prk)
	if err != nil {
		return nil, err
	}
	return EncodeCertificate(der), nil
}

func EncodeCertificate(der []byte) *pem.Block {
	return &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	}
}

func CACertificates(q string) ([]*QueryResult, error) {
	pq := PEMQuery{
		Query:         []string{q, "IsCA: true"},
		CaseSensitive: false,
		Types:         []string{"CERTIFICATE"},
	}
	return pq.QueryPaths(CertPath(), true)
}
