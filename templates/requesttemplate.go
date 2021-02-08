package templates

import (
	"fmt"
	"strings"

	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"gopkg.in/yaml.v3"
)

type RequestTemplate x509.CertificateRequest

func (crt RequestTemplate) String() string {
	return fmt.Sprintf("CSR: %s", crt.Subject.String())
}

func NewRequestTemplate(bl *pem.Block) (*RequestTemplate, error) {
	cr, err := x509.ParseCertificateRequest(bl.Bytes)
	if err != nil {
		return nil, err
	}
	rt := RequestTemplate(*cr)
	return &rt, nil
}
func (cr RequestTemplate) MarshalYAML() (interface{}, error) {
	var pkt *PublicKeyTemplate
	if cr.PublicKey != nil {
		pkt = &PublicKeyTemplate{key: cr.PublicKey}
	}

	return &yamlCSR{
		Version:            cr.Version,
		Subject:            SubjectTemplate(cr.Subject),
		PublicKey:          pkt,
		PublicKeyAlgorithm: PublicKeyAlgorithmTemplate(cr.PublicKeyAlgorithm),
		Signature:          fmt.Sprintf("%x", &cr.Signature),
		SignatureAlgorithm: SignatureAlgorithmTemplate(cr.SignatureAlgorithm),
		Extensions:         ExtensionsTemplateSlice(cr.Extensions),
		ExtraExtensions:    ExtensionsTemplateSlice(cr.ExtraExtensions),

		DNSNames:       strings.Join(cr.DNSNames, ", "),
		EmailAddresses: strings.Join(cr.EmailAddresses, ", "),
		IPAddresses:    cr.IPAddresses,
		URIs:           cr.URIs,
	}, nil
}

func (crt RequestTemplate) UnmarshalYAML(value *yaml.Node) error {
	var cy yamlCSR
	if err := value.Decode(&cy); err != nil {
		return err
	}
	crt.Version = cy.Version
	crt.Subject = pkix.Name(cy.Subject)
	crt.PublicKey = cy.PublicKey
	crt.PublicKeyAlgorithm = x509.PublicKeyAlgorithm(cy.PublicKeyAlgorithm)
	_, _ = fmt.Sscanf(cy.Signature, "%x", &crt.Signature)
	crt.SignatureAlgorithm = x509.SignatureAlgorithm(cy.SignatureAlgorithm)
	crt.Extensions = ExtensionsSlice(cy.Extensions)
	crt.ExtraExtensions = ExtensionsSlice(cy.ExtraExtensions)
	crt.DNSNames = strings.Split(cy.DNSNames, ",")
	crt.EmailAddresses = strings.Split(cy.EmailAddresses, ",")
	crt.IPAddresses = cy.IPAddresses
	crt.URIs = cy.URIs
	return nil
}

// yamlCSR is the yaml version of the csr, with enums tranlated to stings and empty values removed
type yamlCSR struct {
	Version            int                        `yaml:"Version,omitempty"`
	Subject            SubjectTemplate            `yaml:"Subject"`
	PublicKey          *PublicKeyTemplate         `yaml:"PublicKey,omitempty"`
	PublicKeyAlgorithm PublicKeyAlgorithmTemplate `yaml:"PublicKeyAlgorithm,omitempty"`
	Signature          string                     `yaml:"Signature,omitempty"`
	SignatureAlgorithm SignatureAlgorithmTemplate `yaml:"SignatureAlgorithm,omitempty"`
	Extensions         []ExtensionsTemplate       `yaml:"Extensions,omitempty"`
	ExtraExtensions    []ExtensionsTemplate       `yaml:"ExtraExtensions,omitempty"`
	DNSNames           string                     `yaml:"DNSNames,omitempty"`
	EmailAddresses     string                     `yaml:"EmailAddresses,omitempty"`
	IPAddresses        IPAddressTemplate          `yaml:"IPAddresses,omitempty"`
	URIs               URIsTemplate               `yaml:"URIs,omitempty"`
}
