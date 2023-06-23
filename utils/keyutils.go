package utils

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"reflect"
	"strings"
)

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
	if k1Key, ok := k1.(equalKey); ok {
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

const (
	UnknownCurve ECDSACurve = iota
	P224
	P256
	P384
	P521
)

type ECDSACurve int

var ECDSACurveNames = []string{
	"UnknownCurve",
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
	return UnknownCurve
}
