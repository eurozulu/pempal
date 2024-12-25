package utils

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/model"
)

func PublicKeyFromPrivate(prk crypto.PrivateKey) (crypto.PublicKey, error) {
	type i interface {
		Public() crypto.PublicKey
	}
	k, ok := prk.(i)
	if !ok {
		return nil, fmt.Errorf("key type %T not supported", prk)
	}
	return k.Public(), nil
}

func PublicKeyToPem(puk crypto.PublicKey) ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}), err
}

func PrivateKeyToPem(prk crypto.PrivateKey) ([]byte, error) {
	der, err := x509.MarshalPKCS8PrivateKey(prk)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	}), nil
}

func PublicKeyFromPem(pem *pem.Block) (crypto.PublicKey, error) {
	rt := model.ParseResourceTypeFromPEMType(pem.Type)
	if rt == model.UnknownResourceType {
		return nil, fmt.Errorf("unknown resource type")
	}
	switch rt {
	case model.PublicKey:
		return x509.ParsePKIXPublicKey(pem.Bytes)
	case model.PrivateKey:
		prk, err := x509.ParsePKCS8PrivateKey(pem.Bytes)
		if err != nil {
			return nil, err
		}
		return PublicKeyFromPrivate(prk)
	default:
		return nil, fmt.Errorf("invalid key type %q", pem.Type)
	}
}

func PemToPrivateKey(pembytes []byte) (crypto.PrivateKey, error) {
	blk, _ := pem.Decode(pembytes)
	if blk == nil {
		return nil, fmt.Errorf("no pem encoding found")
	}
	return x509.ParsePKCS8PrivateKey(blk.Bytes)
}

func PrivateKeyType(prk crypto.PrivateKey) x509.PublicKeyAlgorithm {
	switch prk.(type) {
	case *rsa.PrivateKey, rsa.PrivateKey:
		return x509.PublicKeyAlgorithm(x509.RSA)
	case *ed25519.PrivateKey, ed25519.PrivateKey:
		return x509.PublicKeyAlgorithm(x509.Ed25519)
	case *ecdsa.PrivateKey, ecdsa.PrivateKey:
		return x509.PublicKeyAlgorithm(x509.ECDSA)
	case *dsa.PrivateKey, dsa.PrivateKey:
		return x509.PublicKeyAlgorithm(x509.DSA)
	default:
		return x509.PublicKeyAlgorithm(x509.UnknownPublicKeyAlgorithm)
	}
}

func PublicKeyType(puk crypto.PublicKey) x509.PublicKeyAlgorithm {
	switch puk.(type) {
	case *rsa.PublicKey, rsa.PublicKey:
		return x509.PublicKeyAlgorithm(x509.RSA)
	case *ed25519.PublicKey, ed25519.PublicKey:
		return x509.PublicKeyAlgorithm(x509.Ed25519)
	case *ecdsa.PublicKey, ecdsa.PublicKey:
		return x509.PublicKeyAlgorithm(x509.ECDSA)
	case *dsa.PublicKey, dsa.PublicKey:
		return x509.PublicKeyAlgorithm(x509.DSA)
	default:
		return x509.PublicKeyAlgorithm(x509.UnknownPublicKeyAlgorithm)
	}
}
