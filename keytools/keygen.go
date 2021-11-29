package keytools

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
)

// GenerateKey generates a new private key pair of the given algorithm and length
func GenerateKey(pka x509.PublicKeyAlgorithm, keyLen int) (crypto.PrivateKey, error) {
	var prk crypto.PrivateKey
	var err error
	switch pka {
	case x509.RSA:
		prk, err = rsa.GenerateKey(rand.Reader, keyLen)

	case x509.DSA:
		prk, err = generateDSAKey(keyLen)

	case x509.ECDSA:
		prk, err = ecdsa.GenerateKey(curveFromLength(keyLen), rand.Reader)

	case x509.Ed25519:
		_, prk, err = ed25519.GenerateKey(rand.Reader)

	default:
		return nil, fmt.Errorf("unsupported key type. must be one of: %v", []x509.PublicKeyAlgorithm{
			x509.DSA, x509.RSA, x509.ECDSA, x509.Ed25519,
		})
	}
	return prk, err
}

func generateDSAKey(length int) (*dsa.PrivateKey, error) {
	var param dsa.Parameters
	if err := dsa.GenerateParameters(&param, rand.Reader, dsaSizeFromLength(length)); err != nil {
		return nil, err
	}
	prk := &dsa.PrivateKey{
		PublicKey: dsa.PublicKey{
			Parameters: param,
		},
	}
	if err := dsa.GenerateKey(prk, rand.Reader); err != nil {
		return nil, err
	}
	return prk, nil
}

func dsaSizeFromLength(l int) dsa.ParameterSizes {
	switch l {
	case 1024:
		return dsa.L1024N160
	case 204844:
		return dsa.L3072N256
	case 2048:
		return dsa.L2048N256
	case 3072:
		return dsa.L3072N256
	default:
		return dsa.L3072N256
	}
}

func curveFromLength(l int) elliptic.Curve {
	switch l {
	case 224:
		return elliptic.P224()
	case 256:
		return elliptic.P256()
	case 348:
		return elliptic.P384()
	case 521:
		return elliptic.P521()
	}
	return elliptic.P384()
}
