package encoders

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"pempal/pemtypes"
	"pempal/templates"
)

type RequestEncoder struct {
}

func (r RequestEncoder) Encode(p *pem.Block) (templates.Template, error) {
	pt := pemtypes.ParsePEMType(p.Type)
	if pt != pemtypes.Request {
		return nil, fmt.Errorf("%s cannot be encoded into a request", p.Type)
	}
	csr, err := x509.ParseCertificateRequest(p.Bytes)
	if err != nil {
		return nil, err
	}
	var t templates.CSRTemplate
	r.ApplyPem(csr, &t)
	return &t, nil
}

func (r RequestEncoder) ApplyPem(csr *x509.CertificateRequest, t *templates.CSRTemplate) {
	t.Version = csr.Version
	t.SignatureAlgorithm = csr.SignatureAlgorithm.String()
	t.PublicKeyAlgorithm = csr.PublicKeyAlgorithm.String()

	t.Subject = &templates.NameTemplate{}
	NameEncoder{}.ApplyPem(&csr.Subject, t.Subject)

	if len(csr.DNSNames) > 0 {
		t.DNSNames = csr.DNSNames
	}
	if len(csr.EmailAddresses) > 0 {
		t.EmailAddresses = csr.EmailAddresses
	}
	if len(csr.IPAddresses) > 0 {
		t.IPAddresses = MarshalIPAddresses(csr.IPAddresses)
	}
	if len(csr.URIs) > 0 {
		t.URIs = MarshalURIs(csr.URIs)
	}

}
