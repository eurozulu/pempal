package key

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
	"pempal/formats/formathelpers"
	"strconv"
)

// MakeKey creates a new key of the given Key Algorithm.
// supported types are:
// RSA
// ECDSA
// Ed25519
// the 'size' is dependent on the key type.
// RSA keys it should be an integer value representing the bit size of the key
// ECDSA it should be the parseCurve size, as one of: "p224", "p256", "p384" or "p521"
// Ed25519 size is ignored.
func MakeKey(keyAlgorithm x509.PublicKeyAlgorithm, size string) ([]byte, error) {
	switch keyAlgorithm {
	case x509.RSA:
		bits, _ := strconv.Atoi(size)
		return makeKeyRSA(bits)

	case x509.ECDSA:
		return makeKeyECDSA(formathelpers.ParseCurve(size))

	case x509.Ed25519:
		return makeKeyEd25519()

	default:
		return nil, fmt.Errorf("%s not a known KeyAlgorithm", keyAlgorithm.String())
	}
}

func makeKeyEd25519() ([]byte, error) {
	pk, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return x509.MarshalPKCS8PrivateKey(pk)
}

func makeKeyECDSA(c elliptic.Curve) ([]byte, error) {
	if c == nil {
		return nil, fmt.Errorf("ECDSA has no parseCurve size set")
	}
	pk, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return nil, err
	}
	return x509.MarshalPKCS8PrivateKey(pk)
}

func makeKeyRSA(bits int) ([]byte, error) {
	if bits < 1 {
		return nil, fmt.Errorf("RSA Key size not stated")
	}
	pk, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return x509.MarshalPKCS8PrivateKey(pk)
}

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
	case *rsa.PublicKey:
		return x509.RSA
	case *ecdsa.PublicKey:
		return x509.ECDSA
	case *ed25519.PublicKey:
		return x509.Ed25519
	case *dsa.PublicKey:
		return x509.DSA
	default:
		return x509.UnknownPublicKeyAlgorithm
	}
}
