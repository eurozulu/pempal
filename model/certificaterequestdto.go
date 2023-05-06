package model

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"fmt"
	"pempal/utils"
)

type CertificateRequestDTO struct {
	Version            int                   `yaml:"version"`
	Signature          string                `yaml:"signature"`
	SignatureAlgorithm string                `yaml:"signature-algorithm"`
	PublicKeyAlgorithm string                `yaml:"public-key-algorithm"`
	PublicKey          *PublicKeyDTO         `yaml:"public-key"`
	Subject            *DistinguishedNameDTO `yaml:"subject"`

	Der          []byte `yaml:"-"`
	Identity     string `yaml:"identity"`
	ResourceType string `yaml:"resource-type"`
}

func (cd CertificateRequestDTO) String() string {
	if cd.PublicKey == nil {
		return ""
	}
	der, err := cd.PublicKey.MarshalBinary()
	if err != nil {
		return ""
	}
	return MD5PublicKey(der)
}

func (cd CertificateRequestDTO) ToCertificateRequest() (*x509.CertificateRequest, error) {
	var puk crypto.PublicKey
	if cd.PublicKey != nil {
		k, err := cd.PublicKey.ToPublicKey()
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key  %v", err)
		}
		puk = k
	}
	var subject pkix.Name
	if cd.Subject != nil {
		subject = cd.Subject.ToName()
	}

	var signature []byte
	if cd.Signature != "" {
		by, err := hex.DecodeString(cd.Signature)
		if err != nil {
			return nil, fmt.Errorf("failed to decodfe signature as hex  %v", err)
		}
		signature = by
	}
	return &x509.CertificateRequest{
		Version:            cd.Version,
		SignatureAlgorithm: utils.ParseSignatureAlgorithm(cd.SignatureAlgorithm),
		PublicKeyAlgorithm: utils.ParsePublicKeyAlgorithm(cd.PublicKeyAlgorithm),
		PublicKey:          puk,
		Subject:            subject,
		Signature:          signature,

		Extensions:      nil,
		ExtraExtensions: nil,
		DNSNames:        nil,
		EmailAddresses:  nil,
		IPAddresses:     nil,
		URIs:            nil,
	}, nil
}

func (crd *CertificateRequestDTO) UnmarshalBinary(data []byte) error {
	csr, err := x509.ParseCertificateRequest(data)
	if err != nil {
		return err
	}
	var pukTemplate *PublicKeyDTO
	if csr.PublicKey != nil {
		pukt, err := NewPublicKeyDTO(csr.PublicKey)
		if err != nil {
			return fmt.Errorf("Failed to parse certificates public key  %v", err)
		} else {
			pukTemplate = &pukt
		}
	}
	crd.Version = csr.Version
	crd.Signature = hex.EncodeToString(csr.Signature)
	crd.SignatureAlgorithm = csr.SignatureAlgorithm.String()
	crd.PublicKeyAlgorithm = csr.PublicKeyAlgorithm.String()
	crd.PublicKey = pukTemplate
	crd.Subject = newDistinguishedNameDTO(csr.Subject)
	crd.Der = csr.Raw
	crd.Identity = crd.String()
	crd.ResourceType = CertificateRequest.String()
	return nil
}

func (cd CertificateRequestDTO) MarshalBinary() (data []byte, err error) {
	if len(cd.Der) == 0 {
		return nil, fmt.Errorf("certificate request is not parsed")
	}
	return cd.Der, nil
}
