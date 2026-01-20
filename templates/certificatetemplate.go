package templates

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"github.com/eurozulu/pempal/model"
	"gopkg.in/yaml.v2"
	"math/big"
	"net"
	"net/url"
	"time"
)

type CertificateTemplate struct {
	Signature          model.Base64Binary       `yaml:"signature,omitempty"`
	SignatureAlgorithm model.SignatureAlgorithm `yaml:"signature-algorithm"`

	PublicKeyAlgorithm model.PublicKeyAlgorithm `yaml:"public-key-algorithm"`
	PublicKey          *model.PublicKey         `yaml:"public-key"`
	Version            int                      `yaml:"version,omitempty"`
	SerialNumber       *model.SerialNumber      `yaml:"serial-number"`
	Issuer             model.DistinguishedName  `yaml:"issuer"`
	Subject            model.DistinguishedName  `yaml:"subject"`
	NotBefore          model.TimeDTO            `yaml:"not-before"`
	NotAfter           model.TimeDTO            `yaml:"not-after"`
	KeyUsage           model.KeyUsage           `yaml:"key-usage,omitempty"`
	ExtKeyUsage        []x509.ExtKeyUsage       `yaml:"ext-key-usage,omitempty"`         // Sequence of extended key usages.
	UnknownExtKeyUsage []asn1.ObjectIdentifier  `yaml:"unknown-ext-key-usage,omitempty"` // Encountered extended key usages unknown to this package.
	SelfSigned         bool                     `yaml:"self-signed,omitempty"`

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

	// UnhandledCriticalExtensions contains a list of extension IDs that
	// were not (fully) processed when parsing. Verify will fail if this
	// slice is non-empty, unless verification is delegated to an OS
	// library which understands all the critical extensions.
	//
	// Users can access these extensions using Extensions and can remove
	// elements from this slice if they believe that they have been
	// handled.
	UnhandledCriticalExtensions []asn1.ObjectIdentifier `yaml:"unhandled-critical-extensions,omitempty"`

	// BasicConstraintsValid indicates whether IsCA, MaxPathLen,
	// and MaxPathLenZero are valid.
	BasicConstraintsValid bool `yaml:"basic-constraints-valid"`
	IsCA                  bool `yaml:"is-ca"`

	// MaxPathLen and MaxPathLenZero indicate the presence and
	// value of the BasicConstraints' "pathLenConstraint".
	//
	// When parsing a certificate, a positive non-zero MaxPathLen
	// means that the field was specified, -1 means it was unset,
	// and MaxPathLenZero being true mean that the field was
	// explicitly set to zero. The case of MaxPathLen==0 with MaxPathLenZero==false
	// should be treated equivalent to -1 (unset).
	//
	// When generating a certificate, an unset pathLenConstraint
	// can be requested with either MaxPathLen == -1 or using the
	// zero value for both MaxPathLen and MaxPathLenZero.
	MaxPathLen int `yaml:"max-path-len"`
	// MaxPathLenZero indicates that BasicConstraintsValid==true
	// and MaxPathLen==0 should be interpreted as an actual
	// maximum path length of zero. Otherwise, that combination is
	// interpreted as MaxPathLen not being set.
	MaxPathLenZero bool `yaml:"max-path-len-zero"`

	SubjectKeyId   model.Base64Binary `yaml:"subject-key-id,omitempty"`
	AuthorityKeyId model.Base64Binary `yaml:"authority-key-id,omitempty"`

	// RFC 5280, 4.2.2.1 (Authority Information Access)
	OCSPServer            []string `yaml:"ocsp-server,omitempty"`
	IssuingCertificateURL []string `yaml:"issuing-certificate-url,omitempty"`

	// Subject Alternate Name Values. (Note that these Values may not be valid
	// if invalid Values were contained within a parsed certificate. For
	// example, an element of DNSNames may not be a valid DNS domain Name.)
	DNSNames       []string   `yaml:"dns-names,omitempty"`
	EmailAddresses []string   `yaml:"email-addresses,omitempty"`
	IPAddresses    []net.IP   `yaml:"ip-addresses,omitempty"`
	URIs           []*url.URL `yaml:"ur-is,omitempty"`

	// Name constraints
	PermittedDNSDomainsCritical bool         `yaml:"permitted-dns-domains-critical,omitempty"` // if true then the Name constraints are marked critical.
	PermittedDNSDomains         []string     `yaml:"permitted-dns-domains,omitempty"`
	ExcludedDNSDomains          []string     `yaml:"excluded-dns-domains,omitempty"`
	PermittedIPRanges           []*net.IPNet `yaml:"permitted-ip-ranges,omitempty"`
	ExcludedIPRanges            []*net.IPNet `yaml:"excluded-ip-ranges,omitempty"`
	PermittedEmailAddresses     []string     `yaml:"permitted-email-addresses,omitempty"`
	ExcludedEmailAddresses      []string     `yaml:"excluded-email-addresses,omitempty"`
	PermittedURIDomains         []string     `yaml:"permitted-uri-domains,omitempty"`
	ExcludedURIDomains          []string     `yaml:"excluded-uri-domains,omitempty"`

	// CRL Distribution Points
	CRLDistributionPoints []string `yaml:"crl-distribution-points,omitempty"`
}

func (c CertificateTemplate) String() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		return ""
	}
	return string(data)
}
func (c *CertificateTemplate) ApplyTo(cert *model.Certificate) {
	if c.Version > 0 {
		cert.Version = c.Version
	}
	if c.SerialNumber != nil {
		cert.SerialNumber = (*big.Int)(c.SerialNumber)
	}

	if len(c.Signature) > 0 {
		cert.Signature = c.Signature
	}
	if x509.SignatureAlgorithm(c.SignatureAlgorithm) != x509.UnknownSignatureAlgorithm {
		cert.SignatureAlgorithm = x509.SignatureAlgorithm(c.SignatureAlgorithm)
	}
	if c.PublicKey != nil {
		cert.PublicKey = c.PublicKey.Public()
		cert.PublicKeyAlgorithm = x509.PublicKeyAlgorithm(model.NewPublicKey(c.PublicKey).PublicKeyAlgorithm())
	}
	if x509.PublicKeyAlgorithm(c.PublicKeyAlgorithm) != x509.UnknownPublicKeyAlgorithm {
		cert.PublicKeyAlgorithm = x509.PublicKeyAlgorithm(c.PublicKeyAlgorithm)
	}

	if !c.Subject.IsEmpty() {
		subject := model.DistinguishedName(cert.Subject)
		subject.Merge(c.Subject)
		cert.Subject = pkix.Name(subject)
	}
	if !c.Issuer.IsEmpty() {
		issuer := model.DistinguishedName(cert.Issuer)
		issuer.Merge(c.Issuer)
		cert.Issuer = pkix.Name(issuer)
	}
	if !time.Time(c.NotBefore).IsZero() {
		cert.NotBefore = time.Time(c.NotBefore)
	}
	if !time.Time(c.NotAfter).IsZero() {
		cert.NotAfter = time.Time(c.NotAfter)
	}
	if c.KeyUsage != 0 {
		cert.KeyUsage = x509.KeyUsage(c.KeyUsage)
	}
	if len(c.Extensions) > 0 {
		cert.Extensions = model.ModelToExtensions(c.Extensions)
	}
	if len(c.ExtraExtensions) > 0 {
		cert.ExtraExtensions = c.ExtraExtensions
	}
	//cert.UnhandledCriticalExtensions = c.UnhandledCriticalExtensions
	if len(c.ExtKeyUsage) > 0 {
		cert.ExtKeyUsage = c.ExtKeyUsage
	}
	if len(c.UnknownExtKeyUsage) > 0 {
		cert.UnknownExtKeyUsage = c.UnknownExtKeyUsage
	}
	cert.IsCA = c.IsCA
	cert.BasicConstraintsValid = c.BasicConstraintsValid
	cert.MaxPathLen = c.MaxPathLen
	cert.MaxPathLenZero = c.MaxPathLenZero
	cert.SubjectKeyId = c.SubjectKeyId
	cert.AuthorityKeyId = c.AuthorityKeyId
	cert.OCSPServer = c.OCSPServer
	cert.IssuingCertificateURL = c.IssuingCertificateURL
	cert.DNSNames = c.DNSNames
	cert.EmailAddresses = c.EmailAddresses
	cert.IPAddresses = c.IPAddresses
	cert.URIs = c.URIs
	cert.PermittedURIDomains = c.PermittedURIDomains
	cert.ExcludedURIDomains = c.ExcludedURIDomains
	cert.PermittedEmailAddresses = c.PermittedEmailAddresses
	cert.ExcludedEmailAddresses = c.ExcludedEmailAddresses
	cert.PermittedDNSDomains = c.PermittedDNSDomains
	cert.ExcludedDNSDomains = c.ExcludedDNSDomains
	cert.PermittedIPRanges = c.PermittedIPRanges
	cert.ExcludedIPRanges = c.ExcludedIPRanges
	cert.CRLDistributionPoints = c.CRLDistributionPoints
}

func NewCertificateTemplate(cert *model.Certificate) *CertificateTemplate {
	return &CertificateTemplate{
		Signature:               cert.Signature,
		SignatureAlgorithm:      model.SignatureAlgorithm(cert.SignatureAlgorithm),
		PublicKeyAlgorithm:      model.PublicKeyAlgorithm(cert.PublicKeyAlgorithm),
		PublicKey:               model.NewPublicKey(cert.PublicKey),
		Version:                 cert.Version,
		SerialNumber:            (*model.SerialNumber)(cert.SerialNumber),
		Issuer:                  model.DistinguishedName(cert.Issuer),
		Subject:                 model.DistinguishedName(cert.Subject),
		NotBefore:               model.TimeDTO(cert.NotBefore),
		NotAfter:                model.TimeDTO(cert.NotAfter),
		KeyUsage:                model.KeyUsage(cert.KeyUsage),
		IsCA:                    cert.IsCA,
		MaxPathLen:              cert.MaxPathLen,
		MaxPathLenZero:          cert.MaxPathLenZero,
		SubjectKeyId:            cert.SubjectKeyId,
		AuthorityKeyId:          cert.AuthorityKeyId,
		OCSPServer:              cert.OCSPServer,
		IssuingCertificateURL:   cert.IssuingCertificateURL,
		DNSNames:                cert.DNSNames,
		EmailAddresses:          cert.EmailAddresses,
		IPAddresses:             cert.IPAddresses,
		URIs:                    cert.URIs,
		PermittedURIDomains:     cert.PermittedURIDomains,
		ExcludedURIDomains:      cert.ExcludedURIDomains,
		PermittedEmailAddresses: cert.PermittedEmailAddresses,
		ExcludedEmailAddresses:  cert.ExcludedEmailAddresses,
		PermittedIPRanges:       cert.PermittedIPRanges,
		ExcludedIPRanges:        cert.ExcludedIPRanges,
		CRLDistributionPoints:   cert.CRLDistributionPoints,
	}
}
