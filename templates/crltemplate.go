package templates

import (
	"crypto/x509/pkix"
)

type CRLTemplate struct {
	Issuer *NameTemplate `yaml:"issuer"`
	// AuthorityKeyId is used to identify the public key associated with the
	// issuing certificate. It is populated from the authorityKeyIdentifier
	// extension when parsing a CRL. It is ignored when creating a CRL; the
	// extension is populated from the issuing certificate itself.
	AuthorityKeyId string `yaml:"authorityKeyId,omitempty"`

	Signature          string `yaml:"signature"`
	SignatureAlgorithm string `yaml:"signature-algorithm"`

	RevokedCertificates []*RevokedCertificateTemplate

	// Number is used to populate the X.509 v2 cRLNumber extension in the CRL,
	// which should be a monotonically increasing sequence number for a given
	// CRL scope and CRL issuer. It is also populated from the cRLNumber
	// extension when parsing a CRL.
	Number int64 `yaml:"number"`

	ThisUpdate string `yaml:"thisUpdate"`
	NextUpdate string `yaml:"nextUpdate"`

	// Extensions contains raw X.509 extensions. When creating a CRL,
	// the Extensions field is ignored, see ExtraExtensions.
	Extensions []pkix.Extension

	// ExtraExtensions contains any additional extensions to add directly to
	// the CRL.
	ExtraExtensions []pkix.Extension
}
