package templates

import (
	"crypto/x509/pkix"
)

type RevokedCertificateTemplate struct {
	SerialNumber   int64            `yaml:"serial-number"`
	RevocationTime string           `yaml:"revocation-time"`
	Extensions     []pkix.Extension `asn1:"optional" yaml:"extensions,omitempty"`
}
