package templates

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"gopkg.in/yaml.v2"
	"net"
	"net/url"
)

type CertificateRequestTemplate struct {
	Version            int                      `yaml:"version,omitempty"`
	Signature          model.Base64Binary       `yaml:"signature,omitempty"`
	SignatureAlgorithm model.SignatureAlgorithm `yaml:"signature-algorithm"`

	PublicKeyAlgorithm model.PublicKeyAlgorithm `yaml:"public-key-algorithm"`
	PublicKey          *model.PublicKey         `yaml:"public-key"`

	Subject model.DistinguishedName `yaml:"subject"`

	// Extensions contains raw X.509 extensions. When parsing certificates,
	// this can be used to extract non-critical extensions that are not
	// parsed by this package. When marshaling certificates, the Extensions
	// field is ignored, see ExtraExtensions.
	Extensions []model.Extension `yaml:"extensions,omitempty"`

	// ExtraExtensions contains extensions to be copied, raw, into any
	// marshaled certificates. Values override any extensions that would
	// otherwise be produced based on the other fields. The ExtraExtensions
	// field is not populated when parsing certificates, see Extensions.
	ExtraExtensions []pkix.Extension `yaml:"extra-extensions,omitempty"`

	// Subject Alternate Name Values. (Note that these Values may not be valid
	// if invalid Values were contained within a parsed certificate. For
	// example, an element of DNSNames may not be a valid DNS domain Name.)
	DNSNames       []string   `yaml:"dns-names,omitempty"`
	EmailAddresses []string   `yaml:"email-addresses,omitempty"`
	IPAddresses    []net.IP   `yaml:"ip-addresses,omitempty"`
	URIs           []*url.URL `yaml:"ur-is,omitempty"`
}

func (C CertificateRequestTemplate) String() string {
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(&C); err != nil {
		logging.Error(err.Error())
		return ""
	}
	return buf.String()
}

func (c *CertificateRequestTemplate) ApplyTo(csr *model.CertificateRequest) {
	if c.Version > 0 {
		csr.Version = c.Version
	}
	if len(c.Signature) > 0 {
		csr.Signature = c.Signature
	}
	if x509.SignatureAlgorithm(c.SignatureAlgorithm) != x509.UnknownSignatureAlgorithm {
		csr.SignatureAlgorithm = x509.SignatureAlgorithm(c.SignatureAlgorithm)
	}
	if c.PublicKey != nil {
		csr.PublicKey = c.PublicKey
		csr.PublicKeyAlgorithm = x509.PublicKeyAlgorithm(model.NewPublicKey(c.PublicKey).PublicKeyAlgorithm())
	}
	if x509.PublicKeyAlgorithm(c.PublicKeyAlgorithm) != x509.UnknownPublicKeyAlgorithm {
		csr.PublicKeyAlgorithm = x509.PublicKeyAlgorithm(c.PublicKeyAlgorithm)
	}

	if !c.Subject.IsEmpty() {
		subject := model.DistinguishedName(csr.Subject)
		subject.Merge(c.Subject)
		csr.Subject = pkix.Name(subject)
	}
	if len(c.Extensions) > 0 {
		csr.Extensions = model.ModelToExtensions(c.Extensions)
	}
	if len(c.ExtraExtensions) > 0 {
		csr.ExtraExtensions = c.ExtraExtensions
	}
	csr.DNSNames = c.DNSNames
	csr.EmailAddresses = c.EmailAddresses
	csr.IPAddresses = c.IPAddresses
	csr.URIs = c.URIs
}

func NewCertificateRequestTemplate(csr *model.CertificateRequest) *CertificateRequestTemplate {
	return &CertificateRequestTemplate{
		Version:            csr.Version,
		Signature:          csr.Signature,
		SignatureAlgorithm: model.SignatureAlgorithm(csr.SignatureAlgorithm),
		PublicKeyAlgorithm: model.PublicKeyAlgorithm(csr.PublicKeyAlgorithm),
		PublicKey:          model.NewPublicKey(csr.PublicKey),
		Subject:            model.DistinguishedName(csr.Subject),
		Extensions:         model.ExtensionsToModel(csr.Extensions),
		ExtraExtensions:    nil,
		DNSNames:           csr.DNSNames,
		EmailAddresses:     csr.EmailAddresses,
		IPAddresses:        csr.IPAddresses,
		URIs:               csr.URIs,
	}
}
