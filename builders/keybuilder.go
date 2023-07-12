package builders

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"strconv"
)

const (
	property_key_algorithm    = "key-algorithm"
	property_key_is_encrypted = "is-encrypted"
	property_key_curve        = "key-curve"
	property_key_length       = "key-length"
)

type keyBuilder struct {
	temps []templates.Template
}

func (k *keyBuilder) Clear() {
	k.temps = nil
}

func (k *keyBuilder) AddTemplate(t ...templates.Template) {
	if len(t) > 0 {
		k.temps = append(k.temps, t...)
	}
}

func (k *keyBuilder) Validate() utils.CompoundErrors {
	return k.validate(k.BuildTemplate())
}

func (k keyBuilder) Build() (resources.Resource, error) {
	t := k.BuildTemplate()
	if errs := k.validate(t); len(errs) > 0 {
		return nil, errs
	}
	kap := t[property_key_algorithm]
	ka := utils.ParsePublicKeyAlgorithm(kap)
	var kl string
	switch ka {
	case x509.RSA:
		kl = t[property_key_length]
	case x509.ECDSA:
		kl = t[property_key_curve]
	default:
		//
	}

	prk, err := createPrivateKey(ka, kl)
	if err != nil {
		return nil, err
	}
	blk, err := utils.PrivateKeyToPEM(prk)
	if err != nil {
		return nil, err
	}

	return resources.NewResource(blk), nil
}

func (k keyBuilder) BuildTemplate() templates.Template {
	return templates.MergeTemplates(k.temps...)
}

func (k keyBuilder) validate(t templates.Template) utils.CompoundErrors {
	var errs utils.CompoundErrors

	var ka x509.PublicKeyAlgorithm
	kap := t[property_key_algorithm]
	if kap == "" {
		errs = append(errs, fmt.Errorf("%s missing", property_key_algorithm))
	} else {
		ka = utils.ParsePublicKeyAlgorithm(kap)
		if ka == x509.UnknownPublicKeyAlgorithm {
			errs = append(errs, fmt.Errorf("%s invalid", property_key_algorithm))
		}
	}
	switch ka {
	case x509.RSA:
		klp := t[property_key_length]
		if klp == "" {
			errs = append(errs, fmt.Errorf("%s missing", property_key_length))
		} else if _, err := strconv.Atoi(klp); err != nil {
			errs = append(errs, fmt.Errorf("%s invalid. Must be an integer number", property_key_length))
		}
	case x509.ECDSA:
		cp := t[property_key_curve]
		if cp == "" {
			errs = append(errs, fmt.Errorf("%s missing", property_key_curve))
		} else if edc := utils.ParseECDSACurve(cp); edc == utils.UnknownCurve {
			errs = append(errs, fmt.Errorf("%s invalid. Must be one of %v", property_key_curve, utils.ECDSACurveNames))
		}
	default:
		// do nothing
	}
	return errs
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
