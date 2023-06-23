package builder

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"strconv"
)

type keyBuilder struct{}

func (kb *keyBuilder) Validate(t templates.Template) []error {
	_, errs := kb.buildDTO(t)
	return errs
}

func (kb *keyBuilder) Build(t templates.Template) (model.Resource, error) {
	dto, errs := kb.buildDTO(t)
	if len(errs) > 0 {
		return nil, CombineErrors(errs)
	}

	pka := utils.ParsePublicKeyAlgorithm(dto.PublicKeyAlgorithm)
	prk, err := createPrivateKey(pka, dto.KeySize)
	if err != nil {
		return nil, err
	}
	der, err := x509.MarshalPKCS8PrivateKey(prk)
	if err != nil {
		return nil, err
	}
	return model.NewResource(&pem.Block{
		Type:  model.PrivateKey.PEMString(),
		Bytes: der,
	}), nil
}

func (kb *keyBuilder) buildDTO(t templates.Template) (*model.PrivateKeyDTO, []error) {
	dto := model.NewDTOForResourceType(model.PrivateKey)
	if err := t.Apply(dto); err != nil {
		return nil, []error{err}
	}
	var errs []error
	keydto := dto.(*model.PrivateKeyDTO)
	pka := utils.ParsePublicKeyAlgorithm(keydto.PublicKeyAlgorithm)
	if pka == x509.UnknownPublicKeyAlgorithm {
		errs = append(errs, fmt.Errorf("public-key-algorithm unknown"))
	} else if err := validateKeyLength(keydto.KeySize, pka); err != nil {
		errs = append(errs, fmt.Errorf("key-size '%s' invalid  %v\n", keydto.KeySize, err))
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return keydto, nil
}

func validateKeyLength(l string, pka x509.PublicKeyAlgorithm) error {
	var err error
	switch pka {
	case x509.RSA:
		_, err = stringToRSAKeyLength(l)

	case x509.ECDSA:
		_, err = stringToCurve(l)

	case x509.Ed25519:
		if l != "" {
			err = fmt.Errorf("unexpected key length '%s' for Ed25519 key. (Key does not support length")
		}
	}
	return err
}

func createPrivateKey(keyAlgorithm x509.PublicKeyAlgorithm, length string) (crypto.PrivateKey, error) {
	switch keyAlgorithm {
	case x509.RSA:
		bits, err := stringToRSAKeyLength(length)
		if err != nil {
			return nil, err
		}
		return rsa.GenerateKey(rand.Reader, bits)

	case x509.ECDSA:
		cv, err := stringToCurve(length)
		if err != nil {
			return nil, err
		}
		return ecdsa.GenerateKey(cv, rand.Reader)

	case x509.Ed25519:
		prk, _, err := ed25519.GenerateKey(rand.Reader)
		return prk, err

	default:
		return nil, fmt.Errorf("%s is not a supported key type", keyAlgorithm.String())
	}
}

func stringToRSAKeyLength(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("no bit length given for rsa key")
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse rsa key bitsize as integer  %v", err)
	}
	return i, nil
}

func stringToCurve(s string) (elliptic.Curve, error) {
	if s == "" {
		return nil, fmt.Errorf("no curve given for ecdsa key")
	}
	cv := utils.ParseECDSACurve(s)
	if cv == utils.UnknownCurve {
		return nil, fmt.Errorf("%s is not a known curve, use one of %v", s[0], utils.ECDSACurveNames)
	}
	return cv.ToCurve(), nil
}
