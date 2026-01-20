package factories

import (
	"crypto/rand"
	"crypto/x509"
	"errors"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/repositories"
	"github.com/eurozulu/pempal/templates"
)

type CertificateRequestFactory struct {
}

func (cf CertificateRequestFactory) Make(t *templates.CertificateRequestTemplate) ([]model.PemResource, error) {
	var newKey *model.PrivateKey
	if t.PublicKey == nil {
		prk, err := CreateDefaultKey()
		if err != nil {
			return nil, err
		}
		newKey = prk
		t.PublicKey = newKey.Public()
	}

	if err := ValidateCertificateRequestTemplate(t); err != nil {
		return nil, err
	}
	csr := &model.CertificateRequest{}
	t.ApplyTo(csr)

	puk := model.NewPublicKey(csr.PublicKey)
	prk, err := repositories.Keys(config.SearchPath()).ByPublicKey(puk)
	if err != nil {
		return nil, err
	}
	der, err := x509.CreateCertificateRequest(rand.Reader, (*x509.CertificateRequest)(csr), prk)
	if err != nil {
		return nil, err
	}
	if err = csr.UnmarshalBinary(der); err != nil {
		return nil, err
	}
	return []model.PemResource{csr}, nil
}

func ValidateCertificateRequestTemplate(t *templates.CertificateRequestTemplate) error {
	if x509.SignatureAlgorithm(t.SignatureAlgorithm) == x509.UnknownSignatureAlgorithm {
		return errors.New("Signature Algorithm is required")
	}
	if t.Subject.String() == "" {
		return errors.New("subject name is required")
	}

	if t.PublicKey == nil {
		return errors.New("public key is required")
	}
	puk := model.NewPublicKey(t.PublicKey)
	if puk.PublicKeyAlgorithm() != t.PublicKeyAlgorithm {
		t.PublicKeyAlgorithm = puk.PublicKeyAlgorithm()
	}
	return nil
}
