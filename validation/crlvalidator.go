package validation

import (
	"fmt"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
)

type CRLValidator struct {
	keyrepo  resources.Keys
	certrepo resources.Certificates
}

func (cv CRLValidator) Validate(template templates.Template) error {
	t, ok := template.(*templates.CRLTemplate)
	if !ok {
		return fmt.Errorf("Invalid CRL template. must be instance of CRLTemplate not %T", template)
	}

	if err := cv.checkIssuerCertificate(t); err != nil {
		return &ValidationError{
			PropertyName: "issuer",
			Message:      err.Error(),
		}
	}
	return nil
}

func (cv CRLValidator) checkIssuerCertificate(t *templates.CRLTemplate) error {
	if t.Issuer.CommonName == "" {
		return fmt.Errorf("Invalid CRL template. issuer missing common name (CN)")
	}
	certs := cv.certrepo.CertificatesBySubject(t.Issuer.ToName())
	if len(certs) == 0 {
		return fmt.Errorf("issuer %q no known", t.Issuer.String())
	}
	if len(certs) > 1 {
		return fmt.Errorf("issuer %q has more than one (%d) certificate", t.Issuer.String(), len(certs))
	}
	return nil
}
