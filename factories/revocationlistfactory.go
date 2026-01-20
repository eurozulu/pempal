package factories

import (
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/repositories"
	"github.com/eurozulu/pempal/templates"
)

type RevocationListFactory struct{}

func (fac RevocationListFactory) Make(ct *templates.RevocationListTemplate) ([]model.PemResource, error) {
	err := ValidateCRLTemplate(ct)
	if err != nil {
		return nil, err
	}

	rlist := &model.RevocationList{}
	if err := ct.ApplyTo(rlist); err != nil {
		return nil, err
	}

	issuer, err := repositories.Certificates(config.SearchPath()).ByName(ct.Issuer)
	if err != nil {
		return nil, fmt.Errorf("invalid issuer: %v", err)
	}

	puk := model.NewPublicKey(issuer.PublicKey)
	prk, err := repositories.Keys(config.SearchPath()).ByPublicKey(puk)
	if err != nil {
		return nil, fmt.Errorf("failed to find key for issuer %s: %v", issuer.Issuer.String(), err)
	}

	der, err := x509.CreateRevocationList(rand.Reader, (*x509.RevocationList)(rlist), (*x509.Certificate)(issuer), prk.Signer())
	if err != nil {
		return nil, err
	}
	if err := rlist.UnmarshalBinary(der); err != nil {
		return nil, err
	}
	return []model.PemResource{rlist}, nil
}

func ValidateCRLTemplate(t templates.Template) error {
	ct, ok := t.(*templates.RevocationListTemplate)
	if !ok {
		return fmt.Errorf("Unexpected template. Expected %T, found %T", ct, t)
	}
	if ct.Issuer.IsEmpty() {
		return fmt.Errorf("requires an issue name")
	}
	if x509.SignatureAlgorithm(ct.SignatureAlgorithm) == x509.UnknownSignatureAlgorithm {
		return fmt.Errorf("invalid signature algorithm")
	}
	return nil
}
