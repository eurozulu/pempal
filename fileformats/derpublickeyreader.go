package fileformats

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"pempal/keytools"
)

type derPublicKeyReader struct{}

func (d derPublicKeyReader) Unmarshal(by []byte) ([]*pem.Block, error) {
	puk, err := ParsePublicKey(by)
	if err != nil {
		return nil, err
	}
	pt := PublicKeyPEMType(keytools.PublicKeyAlgorithm(puk))
	return []*pem.Block{{
		Type:  pt,
		Bytes: by,
	}}, nil
}

func ParsePublicKey(by []byte) (crypto.PublicKey, error) {
	var puk crypto.PublicKey
	var err error
	// try RSA public key
	puk, err = x509.ParsePKCS1PublicKey(by)
	if err == nil && puk != nil {
		return puk, nil
	}
	// Try as a public key
	puk, err = x509.ParsePKIXPublicKey(by)
	if err == nil && puk != nil {
		return puk, nil
	}
	return puk, err
}

// MarshalPublicKey marshals the given key into a pem block.
// The resulting block contains the der encoded bytes of the key and the relevant PEM type.
func MarshalPublicKey(key crypto.PublicKey) (*pem.Block, error) {
	var pemtype string
	var by []byte

	switch vk := key.(type) {
	case *rsa.PublicKey:
		by = x509.MarshalPKCS1PublicKey(vk)
		pemtype = PEM_RSA_PUBLIC_KEY

	case *ecdsa.PublicKey:
		b, err := x509.MarshalPKIXPublicKey(vk)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_EC_PUBLIC_KEY
	case *ed25519.PublicKey:
		b, err := x509.MarshalPKIXPublicKey(vk)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_PUBLIC_KEY

	case *dsa.PrivateKey:
		// TODO: Check this is right.
		k := dsaOpenssl{
			Version: 0,
			P:       vk.P,
			Q:       vk.Q,
			G:       vk.G,
			Pub:     vk.Y,
		}
		b, err := asn1.Marshal(k)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_PUBLIC_KEY
	default:
		return nil, fmt.Errorf("unsupported private key type")
	}
	return &pem.Block{
		Type:  pemtype,
		Bytes: by,
	}, nil
}

func PublicKeyPEMType(ka x509.PublicKeyAlgorithm) string {
	switch ka {
	case x509.RSA:
		return "RSA PUBLIC KEY"
	case x509.DSA:
		return "DSA PUBLIC KEY"
	case x509.ECDSA:
		return "ECDSA PUBLIC KEY"
	default:
		return "PUBLIC KEY"
	}
}
