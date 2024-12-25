package validation

import (
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal/templates"
)

type keyValidator struct {
}

func (k keyValidator) Validate(template templates.Template) error {
	t, ok := template.(*templates.KeyTemplate)
	if !ok {
		return fmt.Errorf("Invalid key template.  Must be a KeyTemplate not a %T", template)
	}

	if t.PublicKeyAlgorithm == 0 {
		return &ValidationError{
			PropertyName: "public-key-algorithm",
			Message:      "unknown or empty PublicKeyAlgorithm",
		}
	}

	if x509.PublicKeyAlgorithm(t.PublicKeyAlgorithm) == x509.RSA {
		if t.KeySize == 0 {
			return &ValidationError{
				PropertyName: "key-size",
				Message:      "RSA keysize can not be zero",
			}
		}
	}
	return nil
}
