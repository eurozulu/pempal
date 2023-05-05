package builders

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"pempal/keymanager"
	"pempal/model"
	"pempal/templates"
)

type CertificateRequestBuilder struct {
	dto  model.CertificateRequestDTO
	keys keymanager.KeyManager
}

func (c CertificateRequestBuilder) ApplyTemplate(tp ...templates.Template) error {
	for _, t := range tp {
		if err := t.Apply(&c.dto); err != nil {
			return err
		}
	}
}

func (c CertificateRequestBuilder) Validate() []error {
	var errs []error
	cert, err := c.dto.ToCertificate()
	if err != nil {
		errs = append(errs, err)
	}
	if cert.PublicKey == nil {
		errs = append(errs, fmt.Errorf("public-key is missing"))
	}
	if cert.Subject.CommonName == "" {
		errs = append(errs, fmt.Errorf("subject.common-name is missing"))
	}
	if cert.Issuer.String() == "" || cert.Issuer.CommonName == "" {
		errs = append(errs, fmt.Errorf("issuerID.common-name is missing"))
	}

}

func (c CertificateRequestBuilder) Build() (resources.Resource, error) {
	cert, err := c.dto.ToCertificate()
	if err != nil {
		return nil, err
	}
	if cert.PublicKey == nil {
		return nil, fmt.Errorf("public-key is missing")
	}
	if cert.Subject.CommonName == "" {
		return nil, fmt.Errorf("subject.common-name is missing")
	}
	if cert.Issuer.String() == "" || cert.Issuer.CommonName == "" {
		return nil, fmt.Errorf("issuerID.common-name is missing")
	}

	der, err := x509.CreateCertificate(rand.Reader, cert, issuerID, puk, prk)
	if err != nil {
		return nil, err
	}
	cr := &resources.CertificateResource{}
	if err = cr.UnmarshalBinary(der); err != nil {
		return nil, err
	}
	return cr, nil
}
