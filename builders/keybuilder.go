package builders

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
	"github.com/eurozulu/pempal/identity"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"strconv"
)

const (
	Property_key_algorithm    = "key-algorithm"
	Property_key_length       = "key-length"
	Property_key_curve        = "key-curve"
	Property_key_is_encrypted = "is-encrypted"
	Property_key_password     = "password"
)

type keyBuilder struct {
}

func (k keyBuilder) Build(t templates.Template) (resources.Resource, error) {
	if errs := k.Validate(t); len(errs) > 0 {
		return nil, errs
	}
	ka := utils.ParsePublicKeyAlgorithm(t[Property_key_algorithm])
	kl := keyAlgorithmLength(ka, t)
	prkPem, err := createPrivateKeyPem(ka, kl)
	if err != nil {
		return nil, err
	}
	if stringToBool(t[Property_key_is_encrypted]) {
		prkPem, err = encryptKey(prkPem, []byte(t[Property_key_password]))
		if err != nil {
			return nil, err
		}
	}
	blk, _ := pem.Decode(prkPem)
	return resources.NewResource(blk), nil
}

func encryptKey(keypem []byte, password []byte) ([]byte, error) {
	nk, err := identity.NewKey("", keypem)
	if err != nil {
		return nil, err
	}
	nk, err = nk.Encrypt(password)
	if err != nil {
		return nil, err
	}
	return []byte(nk.String()), err
}

func (k keyBuilder) Validate(t templates.Template) utils.CompoundErrors {
	var errs utils.CompoundErrors

	var ka x509.PublicKeyAlgorithm
	kas := t[Property_key_algorithm]
	if kas == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_key_algorithm))
	} else {
		ka = utils.ParsePublicKeyAlgorithm(kas)
		if ka == x509.UnknownPublicKeyAlgorithm {
			errs = append(errs, fmt.Errorf("%s invalid", Property_key_algorithm))
		}
	}
	switch ka {
	case x509.RSA:
		klp := t[Property_key_length]
		if klp == "" {
			errs = append(errs, fmt.Errorf("%s missing", Property_key_length))
		} else if _, err := strconv.Atoi(klp); err != nil {
			errs = append(errs, fmt.Errorf("%s invalid. Must be an integer number", Property_key_length))
		}
	case x509.ECDSA:
		cp := t[Property_key_curve]
		if cp == "" {
			errs = append(errs, fmt.Errorf("%s missing", Property_key_curve))
		} else if edc := utils.ParseECDSACurve(cp); edc == utils.UnknownCurve {
			errs = append(errs, fmt.Errorf("%s invalid. Must be one of %v", Property_key_curve, utils.ECDSACurveNames))
		}
	default:
		// do nothing
	}
	var ie bool
	ies, ok := t[Property_key_is_encrypted]
	if !ok || ies == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_key_is_encrypted))
	} else if b, err := strconv.ParseBool(ies); err != nil {
		errs = append(errs, fmt.Errorf("%s invalid, '%s' is not true or false", Property_key_is_encrypted, ie))
	} else {
		ie = b
	}
	pw := t[Property_key_password]
	if ie && pw == "" {
		errs = append(errs, fmt.Errorf("%s missing", Property_key_password))
	}
	return errs
}

func createPrivateKeyPem(keyAlgorithm x509.PublicKeyAlgorithm, length string) ([]byte, error) {
	prk, err := createPrivateKey(keyAlgorithm, length)
	if err != nil {
		return nil, err
	}
	blk, err := utils.PrivateKeyToPEM(prk)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(blk), nil
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

func stringToBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func keyAlgorithmLength(ka x509.PublicKeyAlgorithm, t templates.Template) string {
	switch ka {
	case x509.RSA:
		return t[Property_key_length]
	case x509.ECDSA:
		return t[Property_key_curve]
	default:
		return ""
	}
}
