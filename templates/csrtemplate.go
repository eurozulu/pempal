package templates

import (
	"crypto/x509/pkix"
)

type CSRTemplate struct {
	Version            int           `yaml:"version"`
	Signature          string        `yaml:"signature"`
	SignatureAlgorithm string        `yaml:"signature-algorithm"`
	PublicKeyAlgorithm string        `yaml:"public-key-algorithm"`
	PublicKey          string        `yaml:"public-key"`
	Subject            *NameTemplate `yaml:"subject"`
	Attributes         []pkix.AttributeTypeAndValueSET
	Extensions         []pkix.Extension
	ExtraExtensions    []pkix.Extension
	DNSNames           []string `yaml:"dns-names"`
	EmailAddresses     []string `yaml:"email-addresses"`
	IPAddresses        []string `yaml:"ip-addresses"`
	URIs               []string `yaml:"uris"`
}
