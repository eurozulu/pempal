package revoked

import (
	"crypto/x509/pkix"
	"gopkg.in/yaml.v3"
	"math/big"
	"pempal/formats"
	"pempal/templates"
)

type revokedCertificateFormat struct{}

func (fm revokedCertificateFormat) applyTemplate(cert *pkix.RevokedCertificate, t templates.RevokedCertificateTemplate) error {
	if t.RevocationTime != "" {
		tm, err := formats.parseTime(t.RevocationTime)
		if err != nil {
			return err
		}
		cert.RevocationTime = tm
	}
	if t.SerialNumber > 0 {
		cert.SerialNumber = big.NewInt(t.SerialNumber)
	}
	return nil
}

func (fm revokedCertificateFormat) ensureRevokeCertTemplate(t templates.Template) (*templates.RevokedCertificateTemplate, error) {
	ct, ok := t.(*templates.RevokedCertificateTemplate)
	if !ok {
		ct = &templates.RevokedCertificateTemplate{}
		by, err := yaml.Marshal(t)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(by, ct); err != nil {
			return nil, err
		}
	}
	return ct, nil
}
