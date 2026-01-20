package factories

import (
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/repositories"
	"github.com/eurozulu/pempal/templates"
)

type Factory interface {
	Make(t templates.Template) ([]model.PemResource, error)
}

func Make(t templates.Template) ([]model.PemResource, error) {
	switch tt := t.(type) {
	case *templates.PrivateKeyTemplate:
		return KeyFactory{}.Make(tt)

	case *templates.RevocationListTemplate:
		return RevocationListFactory{}.Make(tt)

	case *templates.CertificateTemplate:
		return CertificateFactory{}.Make(tt)

	case *templates.CertificateRequestTemplate:
		return CertificateRequestFactory{}.Make(tt)

	default:
		return nil, fmt.Errorf("template type: %T is not a known base template", t)
	}
}

func resolveIssuer(dn model.DistinguishedName) (*model.Issuer, error) {
	issuer := repositories.Issuers(config.SearchPath()).MatchByName(dn)
	if len(issuer) == 0 {
		return nil, fmt.Errorf("issuer %s not found", dn)
	}
	if len(issuer) > 1 {
		return nil, fmt.Errorf("issuer %s is ambigious, matches %d issuers", dn, len(issuer))
	}
	return issuer[0], nil
}
