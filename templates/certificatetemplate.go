package templates

import (
	"crypto"
	"crypto/x509"
	"math/big"
	"net"
	"net/url"
	"strings"
	"time"
)

type CertificateTemplate struct {
	Version      int             `yaml:"Version,omitempty"`
	Subject      SubjectTemplate `yaml:"Subject,omitempty"`
	Issuer       SubjectTemplate `yaml:"Issuer,omitempty"`
	SerialNumber *big.Int        `yaml:"SerialNumber,omitempty"`
	NotBefore    time.Time       `yaml:"NotBefore,omitempty"`
	NotAfter     time.Time       `yaml:"NotAfter,omitempty"`
	Fingerprint  string          `yaml:"Fingerprint,omitempty"`

	PublicKey          crypto.PublicKey   `yaml:"PublicKey,omitempty"`
	PublicKeyAlgorithm PublicKeyAlgorithm `yaml:"PublicKeyAlgorithm,omitempty"`
	Signature          []byte             `yaml:"Signature,omitempty"`
	SignatureAlgorithm SignatureAlgorithm `yaml:"SignatureAlgorithm,omitempty"`

	Extensions      []Extension `yaml:"Extensions,omitempty"`
	ExtraExtensions []Extension `yaml:"ExtraExtensions,omitempty"`

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

	FilePath string `yaml:"-"`
	cert     x509.Certificate
}

func (t CertificateTemplate) Location() string {
	return t.FilePath
}

func (t *CertificateTemplate) String() string {
	ca := " "
	if t.IsCA {
		ca = "CA"
	}
	return strings.Join([]string{TemplateType(t), t.Subject.CommonName, ca,
		t.NotAfter.String(), t.Location()}, "\t")
}

func (t CertificateTemplate) CopyCSR(csr *CSRTemplate) {
	t.Subject = csr.Subject
	t.Version = csr.Version
	t.PublicKey = csr.PublicKey
	t.PublicKeyAlgorithm = csr.PublicKeyAlgorithm

	t.Signature = csr.Signature
	t.SignatureAlgorithm = csr.SignatureAlgorithm
	t.Extensions = csr.Extensions
	t.ExtraExtensions = csr.ExtraExtensions

	t.DNSNames = csr.DNSNames
	t.EmailAddresses = csr.EmailAddresses
	t.IPAddresses = csr.IPAddresses
	t.URIs = csr.URIs
	t.FilePath = csr.FilePath
}

// MarshalBinary will marshal the template into an ASN1/der encoded byte block
func (t *CertificateTemplate) MarshalBinary() (data []byte, err error) {
	t.cert.Version = t.Version
	t.cert.SerialNumber = t.SerialNumber
	t.cert.Subject = t.Subject.Subject()
	t.cert.Issuer = t.Issuer.Subject()
	t.cert.NotBefore = t.NotBefore
	t.cert.NotAfter = t.NotAfter

	t.cert.Signature = t.Signature
	t.cert.SignatureAlgorithm = x509.SignatureAlgorithm(t.SignatureAlgorithm)
	t.cert.PublicKey = t.PublicKey
	t.cert.PublicKeyAlgorithm = x509.PublicKeyAlgorithm(t.PublicKeyAlgorithm)

	t.cert.Extensions = ExtensionReslice(t.Extensions)
	t.cert.ExtraExtensions = ExtensionReslice(t.ExtraExtensions)
	t.cert.DNSNames = t.DNSNames
	t.cert.EmailAddresses = t.EmailAddresses
	t.cert.IPAddresses = t.IPAddresses
	t.cert.URIs = t.URIs

	t.cert.KeyUsage = x509.KeyUsage(t.KeyUsage)
	t.cert.ExtKeyUsage = ExtKeyUsagesReslice(t.ExtKeyUsage)
	t.cert.BasicConstraintsValid = t.BasicConstraintsValid
	t.cert.IsCA = t.IsCA
	t.cert.MaxPathLen = t.MaxPathLen
	t.cert.MaxPathLenZero = t.MaxPathLenZero
	t.cert.IssuingCertificateURL = t.IssuingCertificateURL
	t.cert.CRLDistributionPoints = nil
	return t.cert.Raw, nil
}

func (t *CertificateTemplate) UnmarshalBinary(data []byte) error {
	c, err := x509.ParseCertificate(data)
	if err != nil {
		return err
	}
	t.cert = *c

	t.Version = c.Version
	t.SerialNumber = c.SerialNumber
	t.Subject = NewSubjectTemplate(c.Subject)
	t.Issuer = NewSubjectTemplate(c.Issuer)
	t.NotBefore = c.NotBefore
	t.NotAfter = c.NotAfter
	t.Fingerprint = fingerprint(c.Raw)

	t.PublicKeyAlgorithm = PublicKeyAlgorithm(c.PublicKeyAlgorithm)
	t.PublicKey = c.PublicKey

	t.SignatureAlgorithm = SignatureAlgorithm(c.SignatureAlgorithm)
	t.Signature = c.Signature

	t.Extensions = ExtensionSlice(c.Extensions)
	t.ExtraExtensions = ExtensionSlice(c.ExtraExtensions)

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
