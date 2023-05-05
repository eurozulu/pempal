package builders

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"pempal/keymanager"
	"pempal/model"
	"pempal/templates"
)

type CertificateRequestBuilder struct {
	dto  model.CertificateRequestDTO
	keys keymanager.KeyManager
}

func (cb CertificateRequestBuilder) ApplyTemplate(tp ...templates.Template) error {
	for _, t := range tp {
		if err := t.Apply(&cb.dto); err != nil {
			return err
		}
	}
	return nil
}

func (cb CertificateRequestBuilder) Validate() []error {
	var errs []error
	m := cb.RequiredValues()
	for k := range m {
		errs = append(errs, fmt.Errorf("%s invalid", k))
	}
	return errs
}

func (cb CertificateRequestBuilder) RequiredValues() map[string]interface{} {
	m := map[string]interface{}{}
	if cb.dto.Version == 0 {
		m["version"] = 0
	}
	if cb.dto.Signature == "" {
		m["signature"] = 0
	}
	if cb.dto.SignatureAlgorithm == "" {
		m["signature-algorithm"] = x509.UnknownSignatureAlgorithm
	}
	if cb.dto.PublicKeyAlgorithm == "" {
		if cb.dto.PublicKey == nil || cb.dto.PublicKey.PublicKeyAlgorithm == "" {
			m["public-key-algorithm"] = x509.UnknownPublicKeyAlgorithm
		}
	}
	if cb.dto.PublicKey == nil || cb.dto.PublicKey.PublicKey == "" {
		m["public-key"] = nil
	}
	if cb.dto.Subject != nil {
		missing := newDistinguishedNameBuilder(cb.dto.Subject).RequiredValues()
		if len(missing) > 0 {
			m["subject"] = missing
		}
	} else {
		m["subject"] = nil
	}
	return m
}

func (cb CertificateRequestBuilder) Build() (model.PEMResource, error) {
	if errs := cb.Validate(); len(errs) > 0 {
		return nil, fmt.Errorf("%s", collectErrorList(errs, ", "))
	}

	csr, err := cb.dto.ToCertificateRequest()
	if err != nil {
		return nil, err
	}
	puk := csr.PublicKey
	if puk == nil {
		return nil, fmt.Errorf("public-key is missing")
	}
	id, err := keymanager.NewIdentity(puk)
	if err != nil {
		return nil, err
	}

	prk, err := cb.keys.PrivateKey(id)
	if err != nil {
		return nil, fmt.Errorf("failed to locate requester %s private key  %v", id.String(), err)
	}

	der, err := x509.CreateCertificateRequest(rand.Reader, csr, prk)
	if err != nil {
		return nil, err
	}

	return model.NewPemResourceFromBlock(&pem.Block{
		Type:  model.CertificateRequest.PEMString(),
		Bytes: der,
	}), nil
}
