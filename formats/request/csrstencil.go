package request

import (
	"crypto/x509"
	"pempal/formats"
	"pempal/resources"
	"pempal/stencils"
	"pempal/templates"
)

type csrStencil struct {
}

func (st csrStencil) MakeTemplate(r resources.Resource) (templates.Template, error) {
	blk := r.Pem()
	if blk == nil {
		return nil, nil
	}
	csr, err := x509.ParseCertificateRequest(blk.Bytes)
	if err != nil {
		return nil, err
	}

	t := &templates.CSRTemplate{}
	st.copyToTemplate(t, csr)
	return t, nil
}

func (st csrStencil) copyToTemplate(t *templates.CSRTemplate, csr *x509.CertificateRequest) {
	t.Version = csr.Version
	t.SignatureAlgorithm = csr.SignatureAlgorithm.String()
	t.PublicKeyAlgorithm = csr.PublicKeyAlgorithm.String()

	t.Subject = &templates.NameTemplate{}
	stencils.nameStencil{}.copyToTemplate(t.Subject, csr.Subject)

	if len(csr.DNSNames) > 0 {
		t.DNSNames = csr.DNSNames
	}
	if len(csr.EmailAddresses) > 0 {
		t.EmailAddresses = csr.EmailAddresses
	}
	if len(csr.IPAddresses) > 0 {
		t.IPAddresses = formats.marshalIPAddresses(csr.IPAddresses)
	}
	if len(csr.URIs) > 0 {
		t.URIs = formats.marshalURIs(csr.URIs)
	}
}
