package templates

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"golang.org/x/crypto/ssh"
	"log"
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

	Extensions      []pkix.Extension `yaml:"Extensions,omitempty"`
	ExtraExtensions []pkix.Extension `yaml:"ExtraExtensions,omitempty"`

	// Alternate Name values.
	DNSNames       []string   `yaml:"DNSNames,omitempty"`
	EmailAddresses []string   `yaml:"EmailAddresses,omitempty"`
	IPAddresses    []net.IP   `yaml:"IPAddresses,omitempty"`
	URIs           []*url.URL `yaml:"URIs,omitempty"`

	FilePath string `yaml:"-"`
	request  *x509.CertificateRequest
}

func (t CSRTemplate) Location() string {
	return t.FilePath
}

func (t *CSRTemplate) String() string {
	return strings.Join([]string{TemplateType(t), t.Subject.CommonName, t.PublicFingerprint,
		t.Location()}, "\t")
}

func (t CSRTemplate) MarshalBinary() (data []byte, err error) {
	r := t.request
	if r == nil {
		r = &x509.CertificateRequest{}
	}

	r.Version = t.Version
	r.Subject = t.Subject.Subject()
	r.PublicKey = t.PublicKey
	r.PublicKeyAlgorithm = x509.PublicKeyAlgorithm(t.PublicKeyAlgorithm)
	r.Signature = t.Signature
	r.SignatureAlgorithm = x509.SignatureAlgorithm(t.SignatureAlgorithm)
	r.Extensions = t.Extensions
	r.ExtraExtensions = t.ExtraExtensions
	r.DNSNames = t.DNSNames
	r.EmailAddresses = t.EmailAddresses
	r.IPAddresses = t.IPAddresses
	r.URIs = t.URIs
	return r.Raw, nil
}

func (t *CSRTemplate) UnmarshalBinary(by []byte) error {
	csr, err := x509.ParseCertificateRequest(by)
	if err != nil {
		return err
	}
	t.request = csr

	t.Subject = NewSubjectTemplate(csr.Subject)
	t.Version = csr.Version

	spk, err := ssh.NewPublicKey(csr.PublicKey)
	if err != nil {
		log.Println(err)
	} else {
		t.PublicFingerprint = ssh.FingerprintSHA256(spk)
	}

	t.PublicKeyAlgorithm = PublicKeyAlgorithm(csr.PublicKeyAlgorithm)
	t.SignatureAlgorithm = SignatureAlgorithm(csr.SignatureAlgorithm)

	t.Extensions = csr.Extensions
	t.ExtraExtensions = csr.ExtraExtensions

	t.DNSNames = csr.DNSNames
	t.EmailAddresses = csr.EmailAddresses
	t.IPAddresses = csr.IPAddresses
	t.URIs = csr.URIs
	return nil
}
