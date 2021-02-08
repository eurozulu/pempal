package templates

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"gopkg.in/yaml.v3"
)

type CertificateTemplate x509.Certificate

func NewCertificateTemplate(bl *pem.Block) (*CertificateTemplate, error) {
	c, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return nil, err
	}
	ct := CertificateTemplate(*c)
	return &ct, nil
}
func (ct CertificateTemplate) String() string {
	return ct.Subject.CommonName
}

func (ct *CertificateTemplate) UnmarshalYAML(value *yaml.Node) error {
	var yc yamlCertificate
	if err := value.Decode(&yc); err != nil {
		return err
	}

	sn, err := strconv.ParseInt(yc.SerialNumber, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid serial number.  %v", err)
	}
	ct.Version = yc.Version
	ct.SerialNumber = big.NewInt(sn)
	ct.Subject = pkix.Name(yc.Subject)
	ct.SubjectKeyId = []byte(yc.SubjectKeyId)
	ct.Issuer = pkix.Name(yc.Issuer)
	ct.AuthorityKeyId = []byte(yc.AuthorityKeyId)
	ct.NotBefore = yc.NotAfter
	ct.NotAfter = yc.NotAfter
	ct.BasicConstraintsValid = yc.BasicConstraintsValid
	ct.IsCA = yc.IsCA
	ct.MaxPathLen = yc.MaxPathLen
	ct.MaxPathLenZero = yc.MaxPathLenZero
	ct.PublicKey = yc.PublicKey
	ct.PublicKeyAlgorithm = x509.PublicKeyAlgorithm(yc.PublicKeyAlgorithm)
	ct.Signature = []byte(yc.Signature)
	ct.SignatureAlgorithm = x509.SignatureAlgorithm(yc.SignatureAlgorithm)
	ct.KeyUsage = x509.KeyUsage(yc.KeyUsage)
	ct.ExtKeyUsage = ExtKeyUsagesTemplateReslice(yc.ExtKeyUsage)
	ct.Extensions = ExtensionsTemplateReslice(yc.Extensions)
	ct.ExtraExtensions = ExtensionsTemplateReslice(yc.ExtraExtensions)
	ct.DNSNames = strings.Split(yc.DNSNames, ", ")
	ct.EmailAddresses = strings.Split(yc.EmailAddresses, ", ")
	ct.IPAddresses = yc.IPAddresses
	ct.URIs = yc.URIs
	ct.CRLDistributionPoints = strings.Split(yc.CRLDistributionPoints, ", ")
	return nil
}

func (ct CertificateTemplate) MarshalYAML() (interface{}, error) {
	return &yamlCertificate{
		Version:               ct.Version,
		SerialNumber:          fmt.Sprintf("%d", ct.SerialNumber),
		Subject:               SubjectTemplate(ct.Subject),
		SubjectKeyId:          fmt.Sprintf("%x", ct.SubjectKeyId),
		Issuer:                SubjectTemplate(ct.Issuer),
		AuthorityKeyId:        fmt.Sprintf("%x", ct.AuthorityKeyId),
		NotBefore:             ct.NotBefore,
		NotAfter:              ct.NotAfter,
		BasicConstraintsValid: ct.BasicConstraintsValid,
		IsCA:                  ct.IsCA,
		MaxPathLen:            ct.MaxPathLen,
		MaxPathLenZero:        ct.MaxPathLenZero,
		PublicKey:             &PublicKeyTemplate{key: ct.PublicKey},
		PublicKeyAlgorithm:    PublicKeyAlgorithmTemplate(ct.PublicKeyAlgorithm),
		Signature:             fmt.Sprintf("%x", ct.Signature),
		SignatureAlgorithm:    SignatureAlgorithmTemplate(ct.SignatureAlgorithm),
		KeyUsage:              KeyUsageTemplate(ct.KeyUsage),
		ExtKeyUsage:           ExtKeyUsagesTemplateSlice(ct.ExtKeyUsage),
		Extensions:            ExtensionsTemplateSlice(ct.Extensions),
		ExtraExtensions:       ExtensionsTemplateSlice(ct.ExtraExtensions),
		DNSNames:              strings.Join(ct.DNSNames, ", "),
		EmailAddresses:        strings.Join(ct.EmailAddresses, ", "),
		IPAddresses:           ct.IPAddresses,
		URIs:                  ct.URIs,
		CRLDistributionPoints: strings.Join(ct.CRLDistributionPoints, ", "),
	}, nil
}

type yamlCertificate struct {
	Version               int                        `yaml:"Version,omitempty"`
	SerialNumber          string                     `yaml:"SerialNumber,omitempty"`
	Subject               SubjectTemplate            `yaml:"Subject"`
	SubjectKeyId          string                     `yaml:"SubjectKeyId,omitempty"`
	Issuer                SubjectTemplate            `yaml:"Issuer,omitempty"`
	AuthorityKeyId        string                     `yaml:"AuthorityKeyId,omitempty"`
	NotBefore             time.Time                  `yaml:"NotBefore,omitempty"`
	NotAfter              time.Time                  `yaml:"NotAfter,omitempty"`
	BasicConstraintsValid bool                       `yaml:"BasicConstraintsValid,omitempty"`
	IsCA                  bool                       `yaml:"IsCA"`
	MaxPathLen            int                        `yaml:"MaxPathLen,omitempty"`
	MaxPathLenZero        bool                       `yaml:"MaxPathLenZero"`
	PublicKey             *PublicKeyTemplate         `yaml:"PublicKey,omitempty"`
	PublicKeyAlgorithm    PublicKeyAlgorithmTemplate `yaml:"PublicKeyAlgorithm,omitempty"`
	Signature             string                     `yaml:"Signature,omitempty"`
	SignatureAlgorithm    SignatureAlgorithmTemplate `yaml:"SignatureAlgorithm,omitempty"`
	KeyUsage              KeyUsageTemplate           `yaml:"KeyUsage,omitempty"`
	ExtKeyUsage           []ExtKeyUsagesTemplate     `yaml:"ExtKeyUsage,omitempty"`
	Extensions            []ExtensionsTemplate       `yaml:"Extensions,omitempty"`
	ExtraExtensions       []ExtensionsTemplate       `yaml:"ExtraExtensions,omitempty"`
	DNSNames              string                     `yaml:"DNSNames,omitempty"`
	EmailAddresses        string                     `yaml:"EmailAddresses,omitempty"`
	IPAddresses           IPAddressTemplate          `yaml:"IPAddresses,omitempty"`
	URIs                  URIsTemplate               `yaml:"URIs,omitempty"`
	CRLDistributionPoints string                     `yaml:"CRLDistributionPoints,omitempty"`
}
