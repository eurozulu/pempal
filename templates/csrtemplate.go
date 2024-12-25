package templates

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"net"
	"net/url"
	"strings"
)

type CSRTemplate struct {
	Signature          model.Base64Binary       `yaml:"signature,omitempty"`
	SignatureAlgorithm model.SignatureAlgorithm `yaml:"signature-algorithm,omitempty"`

	PublicKeyAlgorithm model.PublicKeyAlgorithm `yaml:"public-key-algorithm,omitempty"`
	PublicKey          model.PublicKeyDTO       `yaml:"public-key,omitempty"`
	ID                 model.KeyId              `yaml:"id"`
	Version            int                      `yaml:"version,omitempty"`
	Subject            model.DistinguishedName  `yaml:"subject"`

	// Extensions contains raw X.509 extensions. When parsing certificates,
	// this can be used to extract non-critical extensions that are not
	// parsed by this package. When marshaling certificates, the Extensions
	// field is ignored, see ExtraExtensions.
	Extensions []pkix.Extension `yaml:"extensions,omitempty"`

	// ExtraExtensions contains extensions to be copied, raw, into any
	// marshaled certificates. Values override any extensions that would
	// otherwise be produced based on the other fields. The ExtraExtensions
	// field is not populated when parsing certificates, see Extensions.
	ExtraExtensions []pkix.Extension `yaml:"extra-extensions,omitempty"`

	// Subject Alternate Name values. (Note that these values may not be valid
	// if invalid values were contained within a parsed certificate. For
	// example, an element of DNSNames may not be a valid DNS domain Name.)
	DNSNames       []string   `yaml:"dns-names,omitempty"`
	EmailAddresses []string   `yaml:"email-addresses,omitempty"`
	IPAddresses    []net.IP   `yaml:"ip-addresses,omitempty"`
	URIs           []*url.URL `yaml:"ur-is,omitempty"`
}

func (ct CSRTemplate) ToCSR() *x509.CertificateRequest {
	return &x509.CertificateRequest{
		Signature:          ct.Signature,
		SignatureAlgorithm: x509.SignatureAlgorithm(ct.SignatureAlgorithm),
		PublicKeyAlgorithm: x509.PublicKeyAlgorithm(ct.PublicKeyAlgorithm),
		PublicKey:          ct.PublicKey,
		Version:            ct.Version,
		Subject:            ct.Subject.ToName(),
		Extensions:         ct.Extensions,
		ExtraExtensions:    ct.ExtraExtensions,
		DNSNames:           ct.DNSNames,
		EmailAddresses:     ct.EmailAddresses,
		IPAddresses:        ct.IPAddresses,
		URIs:               ct.URIs,
	}
}

func (ct CSRTemplate) Name() string {
	return strings.ToLower(model.CertificateRequest.String())
}

func NewCSRTemplate(csr *x509.CertificateRequest) *CSRTemplate {
	ct := CSRTemplate{}
	if csr.PublicKey != nil {
		ct.PublicKey.PublicKey = csr.PublicKey
		if id, err := model.NewKeyIdFromKey(csr.PublicKey); err != nil {
			logging.Error("NewCertificateTemplate", "Failed to read ID from public key  %v", err)
		} else {
			ct.ID = id
		}
	}
	ct.Signature = csr.Signature
	ct.SignatureAlgorithm = model.SignatureAlgorithm(csr.SignatureAlgorithm)
	ct.PublicKeyAlgorithm = model.PublicKeyAlgorithm(csr.PublicKeyAlgorithm)
	ct.Version = csr.Version
	ct.Subject = model.DistinguishedName(csr.Subject)
	ct.Extensions = csr.Extensions
	ct.ExtraExtensions = csr.ExtraExtensions
	ct.DNSNames = csr.DNSNames
	ct.EmailAddresses = csr.EmailAddresses
	ct.IPAddresses = csr.IPAddresses
	ct.URIs = csr.URIs
	return &ct
}
