package validation

import (
	"errors"
	"fmt"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
)

type certificateValidator struct {
	keyrepo  resources.Keys
	certrepo resources.Certificates
}

func (cv certificateValidator) Validate(t templates.Template) error {
	template, ok := t.(*templates.CertificateTemplate)
	if !ok {
		return fmt.Errorf("Invalid certificate template.  Must be a CertificateTemplate not a %T", t)
	}

	if err := cv.checkPublicKey(template); err != nil {
		return &ValidationError{
			PropertyName: "public-key",
			Message:      err.Error(),
		}
	}
	if err := cv.checkSubject(template); err != nil {
		return &ValidationError{
			PropertyName: "subject",
			Message:      err.Error(),
		}
	}
	if err := cv.checkIssuer(template); err != nil {
		return &ValidationError{
			PropertyName: "issuer",
			Message:      err.Error(),
		}
	}
	if err := cv.checkSerialNumber(template); err != nil {
		return &ValidationError{
			PropertyName: "serial-number",
			Message:      err.Error(),
		}
	}
	return nil
}

func (cv certificateValidator) checkPublicKey(t *templates.CertificateTemplate) error {
	if t.PublicKey.PublicKey != nil {
		return nil
	}
	if t.ID != nil {
		puk, err := cv.keyrepo.PublicKeyFromID(t.ID)
		if err != nil {
			return fmt.Errorf("%s is not a valid ID. %v", err)
		}
		t.PublicKey.PublicKey = puk
		return nil
	}
	return errors.New("no 'public-key' or 'id' found")
}

func (cv certificateValidator) checkSubject(t *templates.CertificateTemplate) error {
	if t.Subject.CommonName == "" {
		return fmt.Errorf("invalid certificate template.  missing subject common name (CN)")
	}
	return nil
}

func (cv certificateValidator) checkIssuer(t *templates.CertificateTemplate) error {
	if t.SelfSigned {
		t.Issuer = t.Subject
		return nil
	}
	if t.Issuer.CommonName == "" {
		return fmt.Errorf("Invalid certificate template. issuer missing common name (CN)")
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

func (cv certificateValidator) checkSerialNumber(t *templates.CertificateTemplate) error {
	if t.SerialNumber != nil && len(t.SerialNumber.Bits()) > 0 {
		return nil
	}
	t.SerialNumber = SerialNumberFactory{}.NextSerialNumberFor(t.Issuer.ToName())
	return nil
}
