package factories

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/eurozulu/pempal/config"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/repositories"
	"github.com/eurozulu/pempal/templates"
)

type KeyFactory struct{}

func (kf KeyFactory) Make(kt *templates.PrivateKeyTemplate) ([]model.PemResource, error) {
	err := ValidateKeyTemplate(kt)
	if err != nil {
		return nil, err
	}

	var prk crypto.PrivateKey
	switch x509.PublicKeyAlgorithm(kt.KeyAlgoritum) {
	case x509.RSA:
		prk, err = rsa.GenerateKey(rand.Reader, kt.RSAKeyLength)
	case x509.ECDSA:
		prk, err = ecdsa.GenerateKey(kt.ECDSAKeyCurve.ToCurve(), rand.Reader)
	case x509.Ed25519:
		_, prk, err = ed25519.GenerateKey(rand.Reader)
	case x509.DSA:
		return nil, fmt.Errorf("DSA keys are not supported")
	default:
		return nil, fmt.Errorf("unsupported kt algorithm %s", kt.KeyAlgoritum)
	}
	if err != nil {
		return nil, err
	}
	return []model.PemResource{model.NewPrivateKey(prk)}, nil

}

func ValidateKeyTemplate(kt *templates.PrivateKeyTemplate) error {
	switch x509.PublicKeyAlgorithm(kt.KeyAlgoritum) {
	case x509.Ed25519:

	case x509.RSA:
		if kt.RSAKeyLength <= 0 {
			return errors.New("RSAKeyLength must be greater than zero")
		}
	case x509.ECDSA:
		if kt.ECDSAKeyCurve == model.UnknownCurve {
			return errors.New("key-curve must be set for ECDSA key")
		}
	default:
		return fmt.Errorf("%q is not a supported public key algorithm", kt.KeyAlgoritum)
	}
	return nil
}

func CreateDefaultKey() (*model.PrivateKey, error) {
	temps, err := repositories.Templates(config.TemplatePath()).ExpandedByName(config.DefaultKeyTemplateName())
	if err != nil {
		return nil, fmt.Errorf("The default key template could not be found: %s", err)
	}
	t, err := templates.MergeTemplates(temps)
	res, err := Make(t)
	if err != nil {
		return nil, err
	}
	prk, ok := res[0].(*model.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("The default key template did not create a private key!, found %T", res)
	}
	return prk, nil
}
