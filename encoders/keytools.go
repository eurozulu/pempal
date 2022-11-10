package encoders

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
)

func PublicKeyFromPrivate(prk crypto.PrivateKey) crypto.PublicKey {
	if prk == nil {
		return nil
	}
	switch v := prk.(type) {
	case *rsa.PrivateKey:
		return v.PublicKey
	case *ecdsa.PrivateKey:
		return v.PublicKey
	case *ed25519.PrivateKey:
		return v.Public()
	case *dsa.PrivateKey:
		return v.PublicKey
	default:
		return nil
	}
}
func PublicKeyAlgorithmFromKey(puk crypto.PublicKey) x509.PublicKeyAlgorithm {
	switch puk.(type) {
	case *rsa.PublicKey, rsa.PublicKey:
		return x509.RSA
	case *ecdsa.PublicKey, ecdsa.PublicKey:
		return x509.ECDSA
	case *ed25519.PublicKey, ed25519.PublicKey:
		return x509.Ed25519
	case *dsa.PublicKey, dsa.PublicKey:
		return x509.DSA
	default:
		return x509.UnknownPublicKeyAlgorithm
	}
}

func ParsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	prk, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		prk, err = x509.ParsePKCS1PrivateKey(der)
		if err != nil {
			return nil, err
		}
	}
	return prk, nil
}

func ParsePublicKey(der []byte) (crypto.PublicKey, error) {
	prk, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		prk, err = x509.ParsePKCS1PublicKey(der)
		if err != nil {
			return nil, err
		}
	}
	return prk, nil
}
