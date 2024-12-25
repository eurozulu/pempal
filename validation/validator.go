package validation

import (
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
)

type Validator interface {
	Validate(t templates.Template) error
}

func Validate(t templates.Template) error {
	vdr := ValidatorForTemplate(t.Name())
	if vdr == nil {
		return fmt.Errorf("unknown template: %s type", t.Name())
	}
	return vdr.Validate(t)
}

func ValidatorForTemplate(name string) Validator {
	switch name {
	case "certificate":
		return &certificateValidator{
			certrepo: resources.NewCertificates(config.Config.CertPath),
			keyrepo:  resources.NewKeys(config.Config.CertPath),
		}
	case "privatekey":
		return &keyValidator{}
	case "certificaterequest":
		return &csrValidator{keyrepo: resources.NewKeys(config.Config.CertPath)}
	case "revokationlist":
		return &CRLValidator{
			certrepo: resources.NewCertificates(config.Config.CertPath),
			keyrepo:  resources.NewKeys(config.Config.CertPath),
		}
	default:
		return nil
	}
}
