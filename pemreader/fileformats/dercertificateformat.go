package fileformats

import (
	"crypto/x509"
	"encoding/pem"
)

type derCertificateFormat struct{}
type derCertificateRequestFormat struct{}

func (D derCertificateFormat) Format(by []byte) ([]*pem.Block, error) {
	certs, err := x509.ParseCertificates(by)
	if err != nil {
		return nil, err
	}
	blks := make([]*pem.Block, len(certs))
	for i, c := range certs {
		blks[i] = &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: c.Raw,
		}
	}
	return blks, nil
}

func (D derCertificateRequestFormat) Format(by []byte) ([]*pem.Block, error) {
	csr, err := x509.ParseCertificateRequest(by)
	if err != nil {
		return nil, err
	}
	return []*pem.Block{&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr.Raw,
	}}, nil
}
