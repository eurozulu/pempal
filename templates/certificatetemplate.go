package templates

type CertificateTemplate struct {
	Version               int           `yaml:"version"`
	Signature             string        `yaml:"signature,omitempty"`
	SignatureAlgorithm    string        `yaml:"signature-algorithm"`
	PublicKeyAlgorithm    string        `yaml:"public-key-algorithm"`
	PublicKey             string        `yaml:"public-key"`
	Subject               *NameTemplate `yaml:"subject"`
	Issuer                *NameTemplate `yaml:"issuer"`
	SerialNumber          int64         `yaml:"serial-number"`
	Extensions            []string      `yaml:"extensions"`
	ExtraExtensions       []string      `yaml:"extra-extensions"`
	NotBefore             string        `yaml:"not-before"`
	NotAfter              string        `yaml:"not-after"`
	KeyUsage              string        `yaml:"key-usage,omitempty"`
	ExtKeyUsage           []string      `yaml:"ext-key-usage,omitempty"` // Sequence of extended key usages.
	BasicConstraintsValid bool          `yaml:"basic-constraints-valid"`
	IsCA                  bool          `yaml:"is-ca"`
	MaxPathLen            int           `yaml:"max-path-len"`
	MaxPathLenZero        bool          `yaml:"max-path-len-zero"`
	SubjectKeyId          string        `yaml:"subject-key-id,omitempty"`
	AuthorityKeyId        string        `yaml:"authority-key-id,omitempty"`
	OCSPServer            []string      `yaml:"ocsp-server,omitempty"`
	IssuingCertificateURL []string      `yaml:"issuing-certificate-url,omitempty"`
	DNSNames              []string      `yaml:"dns-names,omitempty"`
	EmailAddresses        []string      `yaml:"email-addresses"`
	IPAddresses           []string      `yaml:"ip-addresses,omitempty"`
	URIs                  []string      `yaml:"uris,omitempty"`

	// Name constraints
	PermittedDNSDomainsCritical bool     `yaml:"permitted-dns-domains-critical,omitempty"`
	PermittedDNSDomains         []string `yaml:"permitted-dns-domains,omitempty"`
	ExcludedDNSDomains          []string `yaml:"excluded-dns-domains,omitempty"`
	PermittedIPRanges           []string `yaml:"permitted-ip-ranges,omitempty"`
	ExcludedIPRanges            []string `yaml:"excluded-ip-ranges,omitempty"`
	PermittedEmailAddresses     []string `yaml:"permitted-email-addresses,omitempty"`
	ExcludedEmailAddresses      []string `yaml:"excluded-email-addresses,omitempty"`
	PermittedURIDomains         []string `yaml:"permitted-uri-domains,omitempty"`
	ExcludedURIDomains          []string `yaml:"excluded-uri-domains,omitempty"`

	// CRL Distribution Points
	CRLDistributionPoints []string
	//PolicyIdentifiers []asn1.ObjectIdentifier
}
