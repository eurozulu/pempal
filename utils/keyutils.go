package utils

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"pempal/logger"
	"reflect"
	"strconv"
	"strings"
)

func CreatePrivateKey(keyAlgorithm x509.PublicKeyAlgorithm, param string) (crypto.PrivateKey, error) {
	switch keyAlgorithm {
	case x509.RSA:
		bits, err := paramToBits(param)
		if err != nil {
			return nil, err
		}
		return rsa.GenerateKey(rand.Reader, bits)

	case x509.ECDSA:
		cv, err := paramToCurve(param)
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

func PublicKeyFromPrivate(key crypto.PrivateKey) (crypto.PublicKey, error) {
	switch v := key.(type) {
	case *rsa.PrivateKey:
		return &v.PublicKey, nil
	case *ed25519.PrivateKey:
		return v.Public(), nil
	case *ecdsa.PrivateKey:
		return &v.PublicKey, nil
	default:
		return nil, fmt.Errorf("%s is an unsupported private key type", reflect.TypeOf(key).Name())
	}
}

func PublicKeyEquals(k1, k2 crypto.PublicKey) bool {
	if k1 == k2 {
		return true
	}

	type equalKey interface {
		Equal(crypto.PublicKey) bool
	}
	if k1Key, ok := k1.(equalKey); !ok {
		logger.Log(logger.Error, "x509: internal error: supported public key does not implement Equal")
	} else {
		return k1Key.Equal(k2)
	}
	return false
}

func PublicKeyAlgorithmFromKey(puk crypto.PublicKey) x509.PublicKeyAlgorithm {
	switch puk.(type) {
	case *rsa.PublicKey:
		return x509.RSA
	case *dsa.PublicKey:
		return x509.DSA
	case *ecdsa.PublicKey:
		return x509.ECDSA
	case *ed25519.PublicKey:
		return x509.Ed25519
	default:
		return x509.UnknownPublicKeyAlgorithm
	}
}

func ParsePublicKeyAlgorithm(s string) x509.PublicKeyAlgorithm {
	for i, pka := range PublicKeyAlgorithms {
		if strings.EqualFold(s, pka) {
			return x509.PublicKeyAlgorithm(i + 1)
		}
	}
	return x509.UnknownPublicKeyAlgorithm
}

var PublicKeyAlgorithms = []string{
	x509.RSA.String(),
	x509.DSA.String(),
	x509.ECDSA.String(),
	x509.Ed25519.String(),
}

func paramToBits(param string) (int, error) {
	if param == "" {
		return 0, fmt.Errorf("no bit length given for rsa key")
	}
	i, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse rsa key bitsize as integer  %v", err)
	}
	return i, nil
}

func paramToCurve(param string) (elliptic.Curve, error) {
	if param == "" {
		return nil, fmt.Errorf("no curve given for ecdsa key")
	}
	cv := ParseECDSACurve(param)
	if cv == Unknown {
		return nil, fmt.Errorf("%s is not a known curve, use one of %v", param[0], ECDSACurveNames)
	}
	return cv.ToCurve(), nil
}

const (
	Unknown ECDSACurve = iota
	P224
	P256
	P384
	P521
)

type ECDSACurve int

var ECDSACurveNames = []string{
	"Unknown",
	"P224",
	"P256",
	"P384",
	"P521",
}

func (curve ECDSACurve) ToCurve() elliptic.Curve {
	switch curve {
	case P224:
		return elliptic.P224()
	case P256:
		return elliptic.P256()
	case P384:
		return elliptic.P384()
	case P521:
		return elliptic.P521()
	default:
		return nil
	}
}

func ParseECDSACurve(s string) ECDSACurve {
	for i, n := range ECDSACurveNames {
		if strings.EqualFold(s, n) {
			return ECDSACurve(i)
		}
	}
	return Unknown
}
