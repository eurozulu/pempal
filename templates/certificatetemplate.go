package templates

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
)

type CertificateTemplate struct {
	Template string          `yaml:"Template"`
	Version  int             `yaml:"Version,omitempty"`
	Subject  SubjectTemplate `yaml:"Subject,omitempty"`
	Issuer   SubjectTemplate `yaml:"Issuer,omitempty"`

	SignatureAlgorithm SignatureAlgorithm `yaml:"SignatureAlgorithm,omitempty"`
	PublicKeyAlgorithm PublicKeyAlgorithm `yaml:"PublicKeyAlgorithm,omitempty"`

	Extensions      []pkix.Extension `yaml:"Extensions,omitempty"`
	ExtraExtensions []pkix.Extension `yaml:"ExtraExtensions,omitempty"`

	// Alternate Name values.
	DNSNames       []string   `yaml:"DNSNames,omitempty"`
	EmailAddresses []string   `yaml:"EmailAddresses,omitempty"`
	IPAddresses    []net.IP   `yaml:"IPAddresses,omitempty"`
	URIs           []*url.URL `yaml:"URIs,omitempty"`

	// certificate only (omitted from requests)
	KeyUsage              KeyUsage      `yaml:"KeyUsage,omitempty"`
	ExtKeyUsage           []ExtKeyUsage `yaml:"ExtKeyUsage,omitempty"`
	BasicConstraintsValid bool          `yaml:"BasicConstraintsValid,omitempty"`
	IsCA                  bool          `yaml:"IsCA,omitempty"`
	MaxPathLen            int           `yaml:"MaxPathLen,omitempty"`
	MaxPathLenZero        bool          `yaml:"MaxPathLenZero,omitempty"`
	IssuingCertificateURL []string      `yaml:"IssuingCertificateURL,omitempty"`
	CRLDistributionPoints []string      `yaml:"CRLDistributionPoints,omitempty"`

	FilePath              string        `yaml:"-"`
	cert                  *x509.Certificate
}

func (t CertificateTemplate) Location() string {
	return t.FilePath
}

func (t *CertificateTemplate) String() string {
	return TemplateString(t)
}

// MarshalBinary will marshal the template into an ASN1/der encoded byte block
func (t *CertificateTemplate) MarshalBinary() (data []byte, err error) {
	//use the unmarshalled cert if available, as the base.
	// otherwise start with an empty cert.
	c := t.cert
	if c == nil {
		c = &x509.Certificate{}
	}
	c.Version = t.Version
	c.Subject = t.Subject.Subject()
	c.Issuer = t.Issuer.Subject()
	c.SignatureAlgorithm = x509.SignatureAlgorithm(t.SignatureAlgorithm)
	c.PublicKeyAlgorithm = x509.PublicKeyAlgorithm(t.PublicKeyAlgorithm)

	c.Extensions = t.Extensions
	c.ExtraExtensions = t.ExtraExtensions
	c.DNSNames = t.DNSNames
	c.EmailAddresses = t.EmailAddresses
	c.IPAddresses = t.IPAddresses
	c.URIs = t.URIs

	c.KeyUsage = x509.KeyUsage(t.KeyUsage)
	c.ExtKeyUsage = ExtKeyUsagesReslice(t.ExtKeyUsage)
	c.BasicConstraintsValid = t.BasicConstraintsValid
	c.IsCA = t.IsCA
	c.MaxPathLen = t.MaxPathLen
	c.MaxPathLenZero = t.MaxPathLenZero
	c.IssuingCertificateURL = t.IssuingCertificateURL
	c.CRLDistributionPoints = nil

	return c.Raw, nil
}

func (t *CertificateTemplate) UnmarshalBinary(data []byte) error {
	c, err := x509.ParseCertificate(data)
	if err != nil {
		return err
	}
	t.cert = c
	t.Template = "CERTIFICATE"

	t.Version = c.Version
	t.Subject = NewSubjectTemplate(c.Subject)
	t.Issuer = NewSubjectTemplate(c.Issuer)

	t.SignatureAlgorithm = SignatureAlgorithm(c.SignatureAlgorithm)
	t.PublicKeyAlgorithm = PublicKeyAlgorithm(c.PublicKeyAlgorithm)

	t.Extensions = c.Extensions
	t.ExtraExtensions = c.ExtraExtensions

	t.DNSNames = c.DNSNames
	t.EmailAddresses = c.EmailAddresses
	t.IPAddresses = c.IPAddresses
	t.URIs = c.URIs

	t.KeyUsage = KeyUsage(c.KeyUsage)
	t.ExtKeyUsage = ExtKeyUsagesSlice(c.ExtKeyUsage)
	t.BasicConstraintsValid = c.BasicConstraintsValid
	t.IsCA = c.IsCA
	t.MaxPathLen = c.MaxPathLen
	t.MaxPathLenZero = c.MaxPathLenZero
	t.IssuingCertificateURL = c.IssuingCertificateURL
	t.CRLDistributionPoints = c.CRLDistributionPoints
	return nil
}

func (t *CertificateTemplate) MarshalYAML() (interface{}, error) {
	t.Template = "CERTIFICATE"
	return t, nil
}

func (t *CertificateTemplate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ct CertificateTemplate
	if err := unmarshal(&ct); err != nil {
		return err
	}
	if ct.Template != "CERTIFICATE" {
		return fmt.Errorf("expected CERTIFICATE template, found '%s'", ct.Template)
	}
	t = &ct
	return nil
}

func (t *CertificateTemplate) MarshalPEM() (*pem.Block, error) {
	by, err := t.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: by,
	}, nil
}

func (t *CertificateTemplate) UnmarshalPEM(bl *pem.Block) error {
	if bl.Type != "CERTIFICATE" {
		return fmt.Errorf("'%s' is not a certificate pem block", bl.Type)
	}
	return t.UnmarshalBinary(bl.Bytes)
}
