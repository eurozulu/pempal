package encoders

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"pempal/pemtypes"
	"pempal/templates"
)

type RequestDecoder struct {
}

func (r RequestDecoder) Decode(t templates.Template) (*pem.Block, error) {
	ct, ok := t.(*templates.CSRTemplate)
	if !ok {
		return nil, fmt.Errorf("template is not a request template")
	}
	var csr x509.CertificateRequest
	r.ApplyTemplate(ct, &csr)
	return &pem.Block{
		Type:  pemtypes.Request.String(),
		Bytes: csr.Raw,
	}, nil
}

func (r RequestDecoder) ApplyTemplate(t *templates.CSRTemplate, csr *x509.CertificateRequest) {
	if t.Version != 0 {
		csr.Version = t.Version
	}
	if t.SignatureAlgorithm != "" {
		csr.SignatureAlgorithm = ParseSignatureAlgorithm(t.SignatureAlgorithm)
	}
	if t.PublicKeyAlgorithm != "" {
		csr.PublicKeyAlgorithm = ParsePublicKeyAlgorithm(t.PublicKeyAlgorithm)
	}
	if t.Subject != nil {
		NameDecoder{}.ApplyTemplate(t.Subject, &csr.Subject)
	}

	//TODO:
	//Attributes         []pkix.AttributeTypeAndValueSET
	//Extensions         []pkix.Extension
	//ExtraExtensions    []pkix.Extension
	if len(t.DNSNames) > 0 {
		csr.DNSNames = t.DNSNames
	}
	if len(t.EmailAddresses) > 0 {
		csr.EmailAddresses = t.EmailAddresses
	}
	if len(t.IPAddresses) > 0 {
		csr.IPAddresses = ParseIPAddresses(t.IPAddresses)
	}
	if len(t.URIs) > 0 {
		csr.URIs = ParseURIs(t.URIs)
	}
}
