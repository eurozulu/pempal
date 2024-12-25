package validation

import (
	"fmt"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
)

type csrValidator struct {
	keyrepo resources.Keys
}

func (cv csrValidator) Validate(template templates.Template) error {
	t, ok := template.(*templates.CSRTemplate)
	if !ok {
		return fmt.Errorf("invaliud CSR template.  Must be a CSRTemplate not a %T", template)
	}
	if err := cv.checkPublicKey(t); err != nil {
		return err
	}
	if err := cv.checkSubject(t); err != nil {
		return &ValidationError{
			PropertyName: "subject",
			Message:      err.Error(),
		}
	}
	return nil
}

func (cv csrValidator) checkSubject(t *templates.CSRTemplate) error {
	if t.Subject.CommonName == "" {
		return fmt.Errorf("invalid CSR template, missing common name (CN)")
	}
	return nil
}

func (cv csrValidator) checkPublicKey(t *templates.CSRTemplate) error {
	if t.PublicKey.PublicKey != nil {
		return nil
	}
	if t.ID != nil {
		puk, err := cv.keyrepo.PublicKeyFromID(t.ID)
		if err != nil {
			return &ValidationError{
				PropertyName: "id",
				Message:      fmt.Sprintf("Invalid ID %q. %v", t.ID, err),
			}
		}
		t.PublicKey.PublicKey = puk
		return nil
	}
	return &ValidationError{
		PropertyName: "public-key",
		Message:      "missing public key or ID",
	}
}
