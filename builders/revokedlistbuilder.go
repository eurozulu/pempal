package builders

import (
	"crypto/x509"
	"fmt"
	"pempal/model"
	"pempal/templates"
	"time"
)

type RevokationListBuilder struct {
	dto model.RevocationListDTO
}

func (rb RevokationListBuilder) ApplyTemplate(tp ...templates.Template) error {
	for _, t := range tp {
		if err := t.Apply(&rb.dto); err != nil {
			return err
		}
	}
	return nil
}

func (rb RevokationListBuilder) Validate() []error {
	var errs []error
	m := rb.RequiredValues()
	for k := range m {
		errs = append(errs, fmt.Errorf("%s invalid", k))
	}
	return errs
}

func (rb RevokationListBuilder) RequiredValues() map[string]interface{} {
	m := map[string]interface{}{}
	if rb.dto.Number == 0 {
		m["number"] = 0
	}
	if rb.dto.Signature == "" {
		m["signature"] = 0
	}
	if rb.dto.SignatureAlgorithm == "" {
		m["signature-algorithm"] = x509.UnknownSignatureAlgorithm
	}
	if rb.dto.ThisUpdate.Before(time.Now()) {
		m["this-update"] = rb.dto.ThisUpdate
	}
	if rb.dto.NextUpdate.Before(time.Now()) || rb.dto.NextUpdate.Before(rb.dto.ThisUpdate) {
		m["next-update"] = rb.dto.NextUpdate
	}
	if len(rb.dto.RevokedCertificates) == 0 {
		m["revoked-certificates"] = []*model.RevokedCertificateDTO{}
	}
	if rb.dto.Issuer != nil {
		missing := newDistinguishedNameBuilder(rb.dto.Issuer).RequiredValues()
		if len(missing) > 0 {
			m["issuer"] = missing
		}
	} else {
		m["issuer"] = nil
	}
	return m
}

func (rb RevokationListBuilder) Build() (model.PEMResource, error) {
	return nil, fmt.Errorf("not yet implemented!")
}
