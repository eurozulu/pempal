package keys

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"strconv"
	"strings"
)

var defaultKeyAlgorithm = x509.RSA
var defaultRSABitLength = 4096
var defaultECDSACurve = elliptic.P384()

// MakeKey creates a new key of the given Key Algorithm.
// supported types are:
// RSA
// ECDSA
// Ed25519
// the 'size' is dependent on the key type.
// RSA keys it should be an integer value representing the bit size of the key
// ECDSA it should be the curve size, as one of: "p224", "p256", "p384" or "p521"
// Ed25519 size is ignored.
func MakeKey(keyAlgorithm x509.PublicKeyAlgorithm, size string) (Key, error) {
	if keyAlgorithm == 0 {
		keyAlgorithm = defaultKeyAlgorithm
	}

	switch keyAlgorithm {
	case x509.RSA:
		bits, _ := strconv.Atoi(size)
		if bits < 1 {
			bits = defaultRSABitLength
		}
		return makeKeyRSA(bits)

	case x509.ECDSA:
		c := curve(size)
		if c == nil {
			c = defaultECDSACurve
		}
		return makeKeyECDSA(c)

	case x509.Ed25519:
		return makeKeyEd25519()

	default:
		return nil, fmt.Errorf("%s not a known KeyAlgorithm", keyAlgorithm.String())
	}
}

func curve(s string) elliptic.Curve {
	switch strings.ToLower(s) {
	case "p224":
		return elliptic.P224()
	case "p256":
		return elliptic.P256()
	case "p384":
		return elliptic.P384()
	case "p521":
		return elliptic.P521()
	default:
		return nil
	}
}

func makeKeyEd25519() (Key, error) {
	pk, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &derKey{pk: pk}, nil
}

func makeKeyECDSA(c elliptic.Curve) (Key, error) {
	pk, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return nil, err
	}
	return &derKey{pk: pk}, nil
}

func makeKeyRSA(bits int) (Key, error) {
	if bits < 1 {
		bits = defaultRSABitLength
	}
	pk, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return &derKey{
		location: "",
		pk:       pk,
	}, nil
}

func ParsePublicKeyAlgorithm(s string) x509.PublicKeyAlgorithm {
	switch s {
	case x509.RSA.String():
		return x509.RSA
	case x509.ECDSA.String():
		return x509.ECDSA
	case x509.Ed25519.String():
		return x509.Ed25519
	default:
		return x509.UnknownPublicKeyAlgorithm
	}
}
