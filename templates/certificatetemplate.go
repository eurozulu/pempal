package templates

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"github.com/eurozulu/pempal/logging"
	"github.com/eurozulu/pempal/model"
	"math/big"
	"net"
	"net/url"
	"strings"
	"time"
)

type CertificateTemplate struct {
	Signature          model.Base64Binary       `yaml:"signature,omitempty"`
	SignatureAlgorithm model.SignatureAlgorithm `yaml:"signature-algorithm"`

	PublicKeyAlgorithm model.PublicKeyAlgorithm `yaml:"public-key-algorithm"`
	PublicKey          model.PublicKeyDTO       `yaml:"public-key"`
	ID                 model.KeyId              `yaml:"id,omitempty"`
	Version            int                      `yaml:"version,omitempty"`
	SerialNumber       *big.Int                 `yaml:"serial-number"`
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

	// Subject Alternate Name values. (Note that these values may not be valid
	// if invalid values were contained within a parsed certificate. For
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

func (ct CertificateTemplate) ToCertificate() *x509.Certificate {
	return &x509.Certificate{
		Signature:                   ct.Signature,
		SignatureAlgorithm:          x509.SignatureAlgorithm(ct.SignatureAlgorithm),
		PublicKeyAlgorithm:          x509.PublicKeyAlgorithm(ct.PublicKeyAlgorithm),
		PublicKey:                   ct.PublicKey.PublicKey,
		Version:                     ct.Version,
		SerialNumber:                ct.SerialNumber,
		Issuer:                      ct.Issuer.ToName(),
		Subject:                     ct.Subject.ToName(),
		NotBefore:                   time.Time(ct.NotBefore),
		NotAfter:                    time.Time(ct.NotAfter),
		KeyUsage:                    x509.KeyUsage(ct.KeyUsage),
		Extensions:                  model.ModelToExtensions(ct.Extensions),
		ExtraExtensions:             ct.ExtraExtensions,
		UnhandledCriticalExtensions: ct.UnhandledCriticalExtensions,
		ExtKeyUsage:                 ct.ExtKeyUsage,
		UnknownExtKeyUsage:          ct.UnknownExtKeyUsage,
		BasicConstraintsValid:       ct.BasicConstraintsValid,
		IsCA:                        ct.IsCA,
		MaxPathLen:                  ct.MaxPathLen,
		MaxPathLenZero:              ct.MaxPathLenZero,
		SubjectKeyId:                ct.SubjectKeyId,
		AuthorityKeyId:              ct.AuthorityKeyId,
		OCSPServer:                  ct.OCSPServer,
		IssuingCertificateURL:       ct.IssuingCertificateURL,
		DNSNames:                    ct.DNSNames,
		EmailAddresses:              ct.EmailAddresses,
		IPAddresses:                 ct.IPAddresses,
		URIs:                        ct.URIs,
		PermittedURIDomains:         ct.PermittedURIDomains,
		ExcludedURIDomains:          ct.ExcludedURIDomains,
		PermittedEmailAddresses:     ct.PermittedEmailAddresses,
		ExcludedEmailAddresses:      ct.ExcludedEmailAddresses,
		PermittedDNSDomains:         ct.PermittedDNSDomains,
		ExcludedDNSDomains:          ct.ExcludedDNSDomains,
		PermittedIPRanges:           ct.PermittedIPRanges,
		ExcludedIPRanges:            ct.ExcludedIPRanges,
		CRLDistributionPoints:       ct.CRLDistributionPoints,
	}
}

func (ct CertificateTemplate) Name() string {
	return strings.ToLower(model.Certificate.String())
}

func NewCertificateTemplate(cert *x509.Certificate) *CertificateTemplate {
	ct := CertificateTemplate{}
	if cert.PublicKey != nil {
		ct.PublicKey.PublicKey = cert.PublicKey
		if id, err := model.NewKeyIdFromKey(cert.PublicKey); err != nil {
			logging.Error("NewCertificateTemplate", "Failed to read ID from public key  %v", err)
		} else {
			ct.ID = id
		}
	}

	ct.Signature = cert.Signature
	ct.SignatureAlgorithm = model.SignatureAlgorithm(cert.SignatureAlgorithm)
	ct.PublicKeyAlgorithm = model.PublicKeyAlgorithm(cert.PublicKeyAlgorithm)
	ct.PublicKey.PublicKey = cert.PublicKey

	ct.Version = cert.Version
	ct.SerialNumber = cert.SerialNumber
	ct.Issuer = model.DistinguishedName(cert.Issuer)
	ct.Subject = model.DistinguishedName(cert.Subject)
	ct.SelfSigned = ct.Issuer.Equals(ct.Subject)
	ct.NotBefore = model.TimeDTO(cert.NotBefore)
	ct.NotAfter = model.TimeDTO(cert.NotAfter)
	ct.KeyUsage = model.KeyUsage(cert.KeyUsage)
	ct.Extensions = model.ExtensionsToModel(cert.Extensions)
	ct.ExtraExtensions = cert.ExtraExtensions
	ct.UnhandledCriticalExtensions = cert.UnhandledCriticalExtensions
	ct.ExtKeyUsage = cert.ExtKeyUsage
	ct.UnknownExtKeyUsage = cert.UnknownExtKeyUsage
	ct.BasicConstraintsValid = cert.BasicConstraintsValid
	ct.IsCA = cert.IsCA
	ct.MaxPathLen = cert.MaxPathLen
	ct.MaxPathLenZero = cert.MaxPathLenZero
	ct.SubjectKeyId = cert.SubjectKeyId
	ct.AuthorityKeyId = cert.AuthorityKeyId
	ct.OCSPServer = cert.OCSPServer
	ct.IssuingCertificateURL = cert.IssuingCertificateURL
	ct.DNSNames = cert.DNSNames
	ct.EmailAddresses = cert.EmailAddresses
	ct.IPAddresses = cert.IPAddresses
	ct.URIs = cert.URIs
	ct.PermittedURIDomains = cert.PermittedURIDomains
	ct.ExcludedURIDomains = cert.ExcludedURIDomains
	ct.PermittedEmailAddresses = cert.PermittedEmailAddresses
	ct.ExcludedEmailAddresses = cert.ExcludedEmailAddresses
	ct.PermittedDNSDomains = cert.PermittedDNSDomains
	ct.ExcludedDNSDomains = cert.ExcludedDNSDomains
	ct.PermittedIPRanges = cert.PermittedIPRanges
	ct.ExcludedIPRanges = cert.ExcludedIPRanges
	ct.CRLDistributionPoints = cert.CRLDistributionPoints
	return &ct
}
