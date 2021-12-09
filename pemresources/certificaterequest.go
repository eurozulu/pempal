package pemresources

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"net"
	"net/url"
	"pempal/fileformats"
	"pempal/keytools"
	"reflect"
)

type CertificateRequest struct {
	PemResource
	Version                     int                     `yaml:"version"`
	Signature                   string                  `yaml:"signature"`
	SignatureAlgorithm          x509.SignatureAlgorithm `yaml:"signature_algorithm"`
	SignatureHash               string                  `yaml:"signature_hash"`
	PublicKeyAlgorithm          x509.PublicKeyAlgorithm `yaml:"public_key_algorithm"`
	PublicKey                   string                  `yaml:"public_key"`
	PublicKeyHash               string                  `yaml:"public_key_hash,omitempty"`
	Subject                     pkix.Name               `yaml:"subject"`
	Extensions                  []pkix.Extension        `yaml:"extensions,omitempty"`
	ExtraExtensions             []pkix.Extension        `yaml:"extra_extensions,omitempty"`
	UnknownExtKeyUsage          []asn1.ObjectIdentifier `yaml:"unknown_ext_key_usage,omitempty"`
	BasicConstraintsValid       bool                    `yaml:"basic_constraints_valid,omitempty"`
	OCSPServer                  []string                `yaml:"ocsp_server,omitempty"`
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
}

func (c CertificateRequest) ResourceId() string {
	return c.PublicKeyHash
}

func (c CertificateRequest) MarshalPem() (*pem.Block, error) {
	blk, err := c.PemResource.MarshalPem()
	if err != nil {
		return nil, err
	}
	csr := x509.CertificateRequest{}

	csr.Version = c.Version

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
		csr.PublicKey = puk
	}
	csr.PublicKeyAlgorithm = c.PublicKeyAlgorithm

	if c.Signature != "" {
		by, err := base64.StdEncoding.DecodeString(c.Signature)
		if err != nil {
			return nil, err
		}
		csr.Signature = by
	}
	csr.SignatureAlgorithm = c.SignatureAlgorithm

	csr.Subject = c.Subject

	csr.Extensions = c.Extensions
	csr.ExtraExtensions = c.ExtraExtensions

	csr.DNSNames = c.DNSNames
	csr.EmailAddresses = c.EmailAddresses
	csr.IPAddresses = c.IPAddresses
	csr.URIs = c.URIs

	by, err := asn1.Marshal(csr)
	if err != nil {
		return nil, err
	}
	blk.Bytes = by
	return blk, nil
}

func (c *CertificateRequest) UnmarshalPem(block *pem.Block) error {
	if err := c.PemResource.UnmarshalPem(block); err != nil {
		return err
	}
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return err
	}
	c.Version = csr.Version

	if csr.PublicKey != nil && !reflect.ValueOf(csr.PublicKey).IsNil() {
		by, err := fileformats.MarshalPublicKey(csr.PublicKey)
		if err != nil {
			return err
		}
		c.PublicKey = base64.StdEncoding.EncodeToString(by.Bytes)
		c.PublicKeyHash = keytools.SHA1HashString(by.Bytes)
	}
	c.PublicKeyAlgorithm = csr.PublicKeyAlgorithm

	c.Signature = base64.StdEncoding.EncodeToString(csr.Signature)
	c.SignatureAlgorithm = csr.SignatureAlgorithm
	c.SignatureHash = keytools.SHA1HashString(csr.Signature)

	c.Subject = csr.Subject

	c.Extensions = csr.Extensions
	c.ExtraExtensions = csr.ExtraExtensions

	c.DNSNames = csr.DNSNames
	c.EmailAddresses = csr.EmailAddresses
	c.IPAddresses = csr.IPAddresses
	c.URIs = csr.URIs
	return nil
}
