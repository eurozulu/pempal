package builders

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"pempal/keymanager"
	"pempal/model"
	"pempal/templates"
	"time"
)

type CertificateBuilder struct {
	dto  model.CertificateDTO
	keys keymanager.KeyManager
}

func (cb CertificateBuilder) ApplyTemplate(tp ...templates.Template) error {
	for _, t := range tp {
		if err := t.Apply(&cb.dto); err != nil {
			return err
		}
	}
	return nil
}

func (cb CertificateBuilder) Validate() []error {
	var errs []error
	m := cb.RequiredValues()
	for k := range m {
		errs = append(errs, fmt.Errorf("%s invalid", k))
	}
	return errs
}

func (cb CertificateBuilder) RequiredValues() map[string]interface{} {
	m := map[string]interface{}{}
	if cb.dto.Version == 0 {
		m["version"] = 0
	}
	if cb.dto.SerialNumber == 0 {
		m["serial-number"] = 0
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
	if cb.dto.Issuer != nil {
		missing := newDistinguishedNameBuilder(cb.dto.Issuer).RequiredValues()
		if len(missing) > 0 {
			m["issuer"] = missing
		}
	} else {
		m["issuer"] = nil
	}
	now := time.Now()
	if cb.dto.NotBefore.Before(now) {
		m["not-before"] = cb.dto.NotBefore
	}
	if cb.dto.NotAfter.Before(cb.dto.NotBefore) || cb.dto.NotAfter.Before(now) {
		m["not-after"] = cb.dto.NotBefore
	}
	return m
}

func (cb CertificateBuilder) Build() (model.PEMResource, error) {
	if errs := cb.Validate(); len(errs) > 0 {
		return nil, fmt.Errorf("%s", collectErrorList(errs, ", "))
	}

	cert, err := cb.dto.ToCertificate()
	if err != nil {
		return nil, err
	}
	issuer, err := cb.keys.User(cert.Issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to locate issuer  %v", err)
	}
	puk := cert.PublicKey
	if puk == nil {
		return nil, fmt.Errorf("public-key is missing")
	}
	der, err := x509.CreateCertificate(rand.Reader, cert, issuer.Certificate(), puk, issuer.Key())
	if err != nil {
		return nil, err
	}

	return model.NewPemResourceFromBlock(&pem.Block{
		Type:  model.Certificate.PEMString(),
		Bytes: der,
	}), nil
}
