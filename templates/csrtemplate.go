package templates

import (
	"crypto"
	"crypto/x509"
	"net"
	"net/url"
	"strings"
)

type CSRTemplate struct {
	Subject SubjectTemplate `yaml:"Subject,omitempty"`
	Version int             `yaml:"Version,omitempty"`

	PublicKey          crypto.PublicKey   `yaml:"PublicKey,omitempty"`
	PublicKeyAlgorithm PublicKeyAlgorithm `yaml:"PublicKeyAlgorithm,omitempty"`
	PublicFingerprint  string             `yaml:"PublicFingerprint,omitempty"`

	Signature          []byte             `yaml:"Signature,omitempty"`
	SignatureAlgorithm SignatureAlgorithm `yaml:"SignatureAlgorithm,omitempty"`

	Extensions      []Extension `yaml:"Extensions,omitempty"`
	ExtraExtensions []Extension `yaml:"ExtraExtensions,omitempty"`

	// Alternate Name values.
	DNSNames       []string   `yaml:"DNSNames,omitempty"`
	EmailAddresses []string   `yaml:"EmailAddresses,omitempty"`
	IPAddresses    []net.IP   `yaml:"IPAddresses,omitempty"`
	URIs           []*url.URL `yaml:"URIs,omitempty"`

	FilePath string `yaml:"-"`
	csr      x509.CertificateRequest
}

func (t CSRTemplate) Location() string {
	return t.FilePath
}

func (t *CSRTemplate) String() string {
	return strings.Join([]string{TemplateType(t), t.Subject.CommonName, t.PublicFingerprint,
		t.Location()}, "\t")
}

func (t CSRTemplate) MarshalBinary() (data []byte, err error) {
	t.csr.Version = t.Version
	t.csr.Subject = t.Subject.Subject()
	t.csr.PublicKey = t.PublicKey
	t.csr.PublicKeyAlgorithm = x509.PublicKeyAlgorithm(t.PublicKeyAlgorithm)
	t.csr.Signature = t.Signature
	t.csr.SignatureAlgorithm = x509.SignatureAlgorithm(t.SignatureAlgorithm)
	t.csr.Extensions = ExtensionReslice(t.Extensions)
	t.csr.ExtraExtensions = ExtensionReslice(t.ExtraExtensions)
	t.csr.DNSNames = t.DNSNames
	t.csr.EmailAddresses = t.EmailAddresses
	t.csr.IPAddresses = t.IPAddresses
	t.csr.URIs = t.URIs
	return t.csr.Raw, nil
}

func (t *CSRTemplate) UnmarshalBinary(by []byte) error {
	csr, err := x509.ParseCertificateRequest(by)
	if err != nil {
		return err
	}
	t.csr = *csr

	t.Subject = NewSubjectTemplate(csr.Subject)
	t.Version = csr.Version

	fpby, err := x509.MarshalPKIXPublicKey(csr.PublicKey)
	if err != nil {
		return err
	}
	t.PublicFingerprint = fingerprint(fpby)

	t.PublicKeyAlgorithm = PublicKeyAlgorithm(csr.PublicKeyAlgorithm)
	t.SignatureAlgorithm = SignatureAlgorithm(csr.SignatureAlgorithm)

	t.Extensions = ExtensionSlice(csr.Extensions)
	t.ExtraExtensions = ExtensionSlice(csr.ExtraExtensions)

	t.DNSNames = csr.DNSNames
	t.EmailAddresses = csr.EmailAddresses
	t.IPAddresses = csr.IPAddresses
	t.URIs = csr.URIs
	return nil
}
