package model

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

type Certificate x509.Certificate

func (c *Certificate) ResourceType() ResourceType {
	return ResourceTypeCertificate
}

func (c *Certificate) String() string {
	return c.Subject.String()
}

func (c *Certificate) Fingerprint() Fingerprint {
	return NewFingerPrint(c.Raw)
}

// MarshalBinary marshalls out certificate in DER format
func (c *Certificate) MarshalBinary() (data []byte, err error) {
	return c.Raw, nil
}

// UnmarshalBinary unmarshalls DER format into a certificate
func (c *Certificate) UnmarshalBinary(data []byte) error {
	cer, err := x509.ParseCertificate(data)
	if err != nil {
		return err
	}
	*c = Certificate(*cer)
	return nil
}

// MarshalText marshals out the certificate in PEM format
func (c *Certificate) MarshalText() (text []byte, err error) {
	der, err := c.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: ResourceTypeCertificate.String(), Bytes: der}), nil
}

// UnmarshalText unmarshals the first PEM certificate in the given pem encoded data
func (c *Certificate) UnmarshalText(text []byte) error {
	var der []byte
	for len(text) > 0 {
		blk, rest := pem.Decode(text)
		if blk == nil {
			break
		}
		if strings.EqualFold(blk.Type, ResourceTypeCertificate.String()) {
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

func NewCertificateFromPem(p *pem.Block) (*Certificate, error) {
	cert := &Certificate{}
	if err := cert.UnmarshalText(pem.EncodeToMemory(p)); err != nil {
		return nil, err
	}
	return cert, nil
}
