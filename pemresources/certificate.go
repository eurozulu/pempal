package pemresources

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"pempal/fileformats"
	"pempal/keytools"
	"reflect"
	"time"
)

type Certificate struct {
	PemResource
	Version                     int                     `yaml:"version"`
	SerialNumber                *big.Int                `yaml:"serial_number"`
	Signature                   string                  `yaml:"signature"`
	SignatureAlgorithm          x509.SignatureAlgorithm `yaml:"signature_algorithm"`
	SignatureHash               string                  `yaml:"signature_hash"`
	PublicKeyAlgorithm          x509.PublicKeyAlgorithm `yaml:"public_key_algorithm"`
	PublicKey                   string                  `yaml:"public_key"`
	PublicKeyHash               string                  `yaml:"public_key_hash,omitempty"`
	Issuer                      pkix.Name               `yaml:"issuer"`
	Subject                     pkix.Name               `yaml:"subject"`
	NotBefore                   time.Time               `yaml:"not_before"`
	NotAfter                    time.Time               `yaml:"not_after"`
	KeyUsage                    x509.KeyUsage           `yaml:"key_usage,omitempty"`
	Extensions                  []pkix.Extension        `yaml:"extensions,omitempty"`
	ExtraExtensions             []pkix.Extension        `yaml:"extra_extensions,omitempty"`
	ExtKeyUsage                 []x509.ExtKeyUsage      `yaml:"ext_key_usage,omitempty"`
	UnknownExtKeyUsage          []asn1.ObjectIdentifier `yaml:"unknown_ext_key_usage,omitempty"`
	BasicConstraintsValid       bool                    `yaml:"basic_constraints_valid,omitempty"`
	IsCA                        bool                    `yaml:"is_ca"`
	MaxPathLen                  int                     `yaml:"max_path_len,omitempty"`
	MaxPathLenZero              bool                    `yaml:"max_path_len_zero"`
	SubjectKeyId                string                  `yaml:"subject_key_id,omitempty"`
	SubjectKeyIdHash            string                  `yaml:"subject_key_id_hash,omitempty"`
	AuthorityKeyId              string                  `yaml:"authority_key_id,omitempty"`
	AuthorityKeyIdHash          string                  `yaml:"authority_key_id_hash,omitempty"`
	OCSPServer                  []string                `yaml:"ocsp_server,omitempty"`
	IssuingCertificateURL       []string                `yaml:"issuing_certificate_url,omitempty"`
	DNSNames                    []string                `yaml:"dns_names,omitempty"`
	EmailAddresses              []string                `yaml:"email_addresses,omitempty"`
	IPAddresses                 []net.IP                `yaml:"ip_addresses,omitempty"`
	URIs                        []*url.URL              `yaml:"ur_is,omitempty"`
	PermittedDNSDomainsCritical bool                    `yaml:"permitted_dns_domains_critical,omitempty"`
	PermittedDNSDomains         []string                `yaml:"permitted_dns_domains,omitempty"`
	ExcludedDNSDomains          []string                `yaml:"excluded_dns_domains,omitempty"`
	PermittedIPRanges           []*net.IPNet            `yaml:"permitted_ip_ranges,omitempty"`
	ExcludedIPRanges            []*net.IPNet            `yaml:"excluded_ip_ranges,omitempty"`
	PermittedEmailAddresses     []string                `yaml:"permitted_email_addresses,omitempty"`
	ExcludedEmailAddresses      []string                `yaml:"excluded_email_addresses,omitempty"`
	PermittedURIDomains         []string                `yaml:"permitted_uri_domains,omitempty"`
	ExcludedURIDomains          []string                `yaml:"excluded_uri_domains,omitempty"`
	CRLDistributionPoints       []string                `yaml:"crl_distribution_points,omitempty"`
}

func (c Certificate) ResourceId() string {
	return c.PublicKeyHash
}

func (c Certificate) MarshalPem() (*pem.Block, error) {
	blk, err := c.PemResource.MarshalPem()
	if err != nil {
		return nil, err
	}
	cert := x509.Certificate{}

	cert.Version = c.Version
	cert.SerialNumber = c.SerialNumber

	if c.PublicKey != "" {
		by, err := base64.StdEncoding.DecodeString(c.PublicKey)
		if err != nil {
			return nil, err
		}
		puk, err := fileformats.ParsePublicKey(by)
		if err != nil {
			return nil, err
		}
		pka := keytools.PublicKeyAlgorithm(puk)
		if c.PublicKeyAlgorithm != pka {
			c.PublicKeyAlgorithm = pka
		}
		cert.PublicKey = puk
	}
	cert.PublicKeyAlgorithm = c.PublicKeyAlgorithm

	if c.Signature != "" {
		by, err := base64.StdEncoding.DecodeString(c.Signature)
		if err != nil {
			return nil, err
		}
		cert.Signature = by
	}
	cert.SignatureAlgorithm = c.SignatureAlgorithm

	if c.SubjectKeyId != "" {
		by, err := base64.StdEncoding.DecodeString(c.SubjectKeyId)
		if err != nil {
			return nil, fmt.Errorf("Invalid subjectkeyId. must be base64 encoded  %w", err)
		}
		cert.SubjectKeyId = by
	}
	if c.AuthorityKeyId != "" {
		by, err := base64.StdEncoding.DecodeString(c.AuthorityKeyId)
		if err != nil {
			return nil, fmt.Errorf("Invalid AuthorityKeyId. must be base64 encoded  %w", err)
		}
		cert.AuthorityKeyId = by
	}

	cert.Issuer = c.Issuer
	cert.Subject = c.Subject

	cert.NotBefore = c.NotBefore
	cert.NotAfter = c.NotAfter
	cert.KeyUsage = c.KeyUsage
	cert.Extensions = c.Extensions
	cert.ExtraExtensions = c.ExtraExtensions
	cert.ExtKeyUsage = c.ExtKeyUsage
	cert.IsCA = c.IsCA
	cert.MaxPathLen = c.MaxPathLen
	cert.MaxPathLenZero = c.MaxPathLenZero

	cert.IssuingCertificateURL = c.IssuingCertificateURL
	cert.DNSNames = c.DNSNames
	cert.EmailAddresses = c.EmailAddresses
	cert.IPAddresses = c.IPAddresses
	cert.URIs = c.URIs
	cert.CRLDistributionPoints = c.CRLDistributionPoints

	by, err := asn1.Marshal(cert)
	if err != nil {
		return nil, err
	}
	blk.Bytes = by
	return blk, nil
}

func (c *Certificate) UnmarshalPem(block *pem.Block) error {
	if c.PemType == "" {
		c.PemResource.PemType = fileformats.PEM_CERTIFICATE
	}
	if err := c.PemResource.UnmarshalPem(block); err != nil {
		return err
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	c.Version = cert.Version
	c.SerialNumber = cert.SerialNumber

	if cert.PublicKey != nil && !reflect.ValueOf(cert.PublicKey).IsNil() {
		by, err := fileformats.MarshalPublicKey(cert.PublicKey)
		if err != nil {
			return err
		}
		c.PublicKey = base64.StdEncoding.EncodeToString(by.Bytes)
		c.PublicKeyHash = keytools.SHA1HashString(by.Bytes)
	}
	c.PublicKeyAlgorithm = cert.PublicKeyAlgorithm

	c.Signature = base64.StdEncoding.EncodeToString(cert.Signature)
	c.SignatureAlgorithm = cert.SignatureAlgorithm
	c.SignatureHash = keytools.SHA1HashString(cert.Signature)

	c.SubjectKeyId = base64.StdEncoding.EncodeToString(cert.SubjectKeyId)
	c.SubjectKeyIdHash = keytools.SHA1HashString(cert.SubjectKeyId)
	c.AuthorityKeyId = base64.StdEncoding.EncodeToString(cert.AuthorityKeyId)
	c.AuthorityKeyIdHash = keytools.SHA1HashString(cert.AuthorityKeyId)

	c.Issuer = cert.Issuer
	c.Subject = cert.Subject

	c.NotBefore = cert.NotBefore
	c.NotAfter = cert.NotAfter
	c.KeyUsage = cert.KeyUsage
	c.Extensions = cert.Extensions
	c.ExtraExtensions = cert.ExtraExtensions
	c.ExtKeyUsage = cert.ExtKeyUsage
	c.IsCA = cert.IsCA
	c.MaxPathLen = cert.MaxPathLen
	c.MaxPathLenZero = cert.MaxPathLenZero

	c.IssuingCertificateURL = cert.IssuingCertificateURL
	c.DNSNames = cert.DNSNames
	c.EmailAddresses = cert.EmailAddresses
	c.IPAddresses = cert.IPAddresses
	c.URIs = cert.URIs
	c.CRLDistributionPoints = cert.CRLDistributionPoints
	return nil
}
