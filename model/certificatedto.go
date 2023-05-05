package model

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"fmt"
	"pempal/utils"
	"time"
)

type CertificateDTO struct {
	Version               int                   `yaml:"version"`
	SerialNumber          SerialNumber          `yaml:"serial-number"`
	Signature             string                `yaml:"signature"`
	SignatureAlgorithm    string                `yaml:"signature-algorithm"`
	PublicKeyAlgorithm    string                `yaml:"public-key-algorithm"`
	PublicKey             *PublicKeyDTO         `yaml:"public-key"`
	Issuer                *DistinguishedNameDTO `yaml:"issuer"`
	Subject               *DistinguishedNameDTO `yaml:"subject"`
	NotBefore             time.Time             `yaml:"not-before"`
	NotAfter              time.Time             `yaml:"not-after"`
	IsCA                  bool                  `yaml:"is-certauth,omitempty"`
	BasicConstraintsValid bool                  `yaml:"basic-constraints-valid,omitempty"`
	MaxPathLen            int                   `yaml:"max-path-len,omitempty"`
	MaxPathLenZero        bool                  `yaml:"max-path-len-zero,omitempty"`

	Der      []byte `yaml:"-"`
	Identity string `yaml:"identity"`
}

func (cd CertificateDTO) String() string {
	if cd.PublicKey == nil {
		return ""
	}
	der, err := cd.PublicKey.MarshalBinary()
	if err != nil {
		return ""
	}
	return MD5PublicKey(der)
}

func (cd CertificateDTO) ToCertificate() (*x509.Certificate, error) {
	var puk crypto.PublicKey
	if cd.PublicKey.PublicKey != "" {
		k, err := cd.PublicKey.ToPublicKey()
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key  %v", err)
		}
		puk = k
	}
	var issuer, subject pkix.Name
	if cd.Issuer != nil {
		issuer = cd.Issuer.ToName()
	}
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
	return &x509.Certificate{
		Version:            cd.Version,
		SerialNumber:       cd.SerialNumber.ToBigInt(),
		SignatureAlgorithm: utils.ParseSignatureAlgorithm(cd.SignatureAlgorithm),
		PublicKeyAlgorithm: utils.ParsePublicKeyAlgorithm(cd.PublicKeyAlgorithm),
		PublicKey:          puk,
		Issuer:             issuer,
		Subject:            subject,
		NotBefore:          cd.NotBefore,
		NotAfter:           cd.NotAfter,
		Signature:          signature,

		KeyUsage:                    0,
		Extensions:                  nil,
		ExtraExtensions:             nil,
		UnhandledCriticalExtensions: nil,
		ExtKeyUsage:                 nil,
		UnknownExtKeyUsage:          nil,

		BasicConstraintsValid:       cd.BasicConstraintsValid,
		IsCA:                        cd.IsCA,
		MaxPathLen:                  cd.MaxPathLen,
		MaxPathLenZero:              cd.MaxPathLenZero,
		SubjectKeyId:                nil,
		AuthorityKeyId:              nil,
		OCSPServer:                  nil,
		IssuingCertificateURL:       nil,
		DNSNames:                    nil,
		EmailAddresses:              nil,
		IPAddresses:                 nil,
		URIs:                        nil,
		PermittedDNSDomainsCritical: false,
		PermittedDNSDomains:         nil,
		ExcludedDNSDomains:          nil,
		PermittedIPRanges:           nil,
		ExcludedIPRanges:            nil,
		PermittedEmailAddresses:     nil,
		ExcludedEmailAddresses:      nil,
		PermittedURIDomains:         nil,
		ExcludedURIDomains:          nil,
		CRLDistributionPoints:       nil,
		PolicyIdentifiers:           nil,
	}, nil
}

func (cd CertificateDTO) MarshalBinary() (data []byte, err error) {
	if len(cd.Der) == 0 {
		return nil, fmt.Errorf("certificate is not parsed")
	}
	return cd.Der, nil
}

func (cd *CertificateDTO) UnmarshalBinary(data []byte) error {
	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return err
	}
	var sn SerialNumber
	if cert.SerialNumber != nil {
		sn = SerialNumber(cert.SerialNumber.Uint64())
	}

	var pukTemplate *PublicKeyDTO
	if cert.PublicKey != nil {
		pukt, err := NewPublicKeyDTO(cert.PublicKey)
		if err != nil {
			return fmt.Errorf("Failed to parse certificates public key  %v", err)
		} else {
			pukTemplate = &pukt
		}
	}
	cd.Version = cert.Version
	cd.SerialNumber = sn
	cd.Signature = hex.EncodeToString(cert.Signature)
	cd.SignatureAlgorithm = cert.SignatureAlgorithm.String()
	cd.PublicKeyAlgorithm = cert.PublicKeyAlgorithm.String()
	cd.PublicKey = pukTemplate
	cd.Issuer = newDistinguishedNameDTO(cert.Issuer)
	cd.Subject = newDistinguishedNameDTO(cert.Subject)
	cd.NotBefore = cert.NotBefore
	cd.NotAfter = cert.NotAfter
	cd.IsCA = cert.IsCA
	cd.BasicConstraintsValid = cert.BasicConstraintsValid
	cd.MaxPathLen = cert.MaxPathLen
	cd.MaxPathLenZero = cert.MaxPathLenZero

	cd.Der = cert.Raw
	cd.Identity = cd.String()
	return nil
}
