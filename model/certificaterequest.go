package model

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

type CertificateRequest x509.CertificateRequest

func (c *CertificateRequest) ResourceType() ResourceType {
	return ResourceTypeCertificateRequest
}

func (c *CertificateRequest) String() string {
	return c.Subject.String()
}

func (c *CertificateRequest) Fingerprint() Fingerprint {
	return NewFingerPrint(c.Raw)
}

func (c *CertificateRequest) MarshalBinary() (data []byte, err error) {
	return c.Raw, nil
}

func (c *CertificateRequest) UnmarshalBinary(data []byte) error {
	cer, err := x509.ParseCertificateRequest(data)
	if err != nil {
		return err
	}
	*c = CertificateRequest(*cer)
	return nil
}

func (c *CertificateRequest) MarshalText() (text []byte, err error) {
	der, err := c.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: ResourceTypeCertificateRequest.String(), Bytes: der}), nil
}

func (c *CertificateRequest) UnmarshalText(text []byte) error {
	var der []byte
	for len(text) > 0 {
		blk, rest := pem.Decode(text)
		if blk == nil {
			break
		}
		if strings.EqualFold(blk.Type, ResourceTypeCertificateRequest.String()) {
			der = blk.Bytes
			break
		}
		text = rest
	}
	if der == nil {
		return errors.New("no certificate PEM found")
	}
	return c.UnmarshalBinary(der)
}

func NewCertificateRequestFromPem(blk *pem.Block) (*CertificateRequest, error) {
	csr := &CertificateRequest{}
	if err := csr.UnmarshalText(pem.EncodeToMemory(blk)); err != nil {
		return nil, err
	}
	return csr, nil
}
