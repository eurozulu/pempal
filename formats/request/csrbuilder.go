package request

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"pempal/formats/dname"
	"pempal/formats/formathelpers"
	"pempal/resources"
	"pempal/templates"
)

type csrBuilder struct {
	csrTemp  templates.CSRTemplate
	location string

	prk crypto.PrivateKey
	puk crypto.PublicKey
}

func (fm *csrBuilder) SetLocation(l string) {
	fm.location = l
}

func (fm *csrBuilder) AddTemplate(ts ...templates.Template) error {
	for _, t := range ts {
		by, err := yaml.Marshal(t)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(by, &fm.csrTemp); err != nil {
			return err
		}
	}
	return nil
}

func (fm csrBuilder) Template() templates.Template {
	return &fm.csrTemp
}

func (fm csrBuilder) Build() (resources.Resources, error) {
	if err := fm.validateTemplate(fm.csrTemp); err != nil {
		return nil, err
	}

	prk, puk, err := formathelpers.ResolveKeys(fm.csrTemp.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve private key %v", err)
	}

	var request x509.CertificateRequest
	fm.applyTemplate(&request, fm.csrTemp)
	if request.Subject.CommonName == "" {
		return nil, fmt.Errorf("missing subject.common-name")
	}
	request.PublicKey = puk

	by, err := x509.CreateCertificateRequest(rand.Reader, &request, prk)
	if err != nil {
		return nil, err
	}
	return resources.Resources{resources.NewResource("", &pem.Block{
		Type:  resources.Request.String(),
		Bytes: by,
	})}, nil
}

func (fm csrBuilder) validateTemplate(ct templates.CSRTemplate) error {
	if ct.Subject == nil {
		return fmt.Errorf("missing subject")
	}
	if ct.Subject.CommonName == "" {
		return fmt.Errorf("missing subject.common-name")
	}
	return nil
}

func (fm csrBuilder) applyTemplate(csr *x509.CertificateRequest, t templates.CSRTemplate) {
	if t.Version != 0 {
		csr.Version = t.Version
	}
	if t.SignatureAlgorithm != "" {
		csr.SignatureAlgorithm = formathelpers.ParseSignatureAlgorithm(t.SignatureAlgorithm)
	}
	if t.PublicKeyAlgorithm != "" {
		csr.PublicKeyAlgorithm = formathelpers.ParsePublicKeyAlgorithm(t.PublicKeyAlgorithm)
	}
	if t.Subject != nil {
		dname.NewNameBuilder(t.Subject).ApplyTemplate(&csr.Subject)
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
		csr.IPAddresses = formathelpers.ParseIPAddresses(t.IPAddresses)
	}
	if len(t.URIs) > 0 {
		csr.URIs = formathelpers.ParseURIs(t.URIs)
	}
}

func NewRequest() *csrBuilder {
	return &csrBuilder{}
}
