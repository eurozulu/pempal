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
	"math/big"
	"pempal/keytools"
)

type derPrivateKeyReader struct{}

func (d derPrivateKeyReader) Unmarshal(by []byte) ([]*pem.Block, error) {
	prk, err := ParsePrivateKey(by)
	if err != nil {
		return nil, err
	}

	pka := keytools.PublicKeyAlgorithmFromPrivate(prk)
	pt := PrivateKeyPEMType(pka)
	blocks := []*pem.Block{{
		Type:  pt,
		Bytes: by,
	}}

	// see if public key is available (unencryoted)
	puk := keytools.PublicKeyFromPrivate(prk)
	if puk != nil {
		pby, err := x509.MarshalPKIXPublicKey(puk)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, &pem.Block{
			Type:  PublicKeyPEMType(keytools.PublicKeyAlgorithm(puk)),
			Bytes: pby,
		})
	}
	return blocks, nil
}

func PrivateKeyPEMType(ka x509.PublicKeyAlgorithm) string {
	switch ka {
	case x509.RSA:
		return "RSA PRIVATE KEY"
	case x509.DSA:
		return "DSA PRIVATE KEY"
	case x509.ECDSA:
		return "EC PRIVATE KEY"
	default: // x509.Ed25519
		return "PRIVATE KEY"
	}
}
func ParsePrivateKey(by []byte) (crypto.PrivateKey, error) {
	var prk crypto.PrivateKey
	var err error

	// attempt to parse as general private key
	prk, err = x509.ParsePKCS8PrivateKey(by)
	if err == nil && prk != nil {
		return prk, nil
	}
	// try as rsa private
	prk, err = x509.ParsePKCS1PrivateKey(by)
	if err == nil && prk != nil {
		return prk, nil
	}
	// try as EC private
	prk, err = x509.ParseECPrivateKey(by)
	if err == nil && prk != nil {
		return prk, nil
	}
	return prk, err
}

// MarshalPrivateKey marshals the given key into a pem block.
// The resulting block contains the der encoded bytes of the key and the relevant PEM type.
func MarshalPrivateKey(key crypto.PrivateKey) (*pem.Block, error) {
	var pemtype string
	var by []byte

	switch vk := key.(type) {
	case *rsa.PrivateKey:
		by = x509.MarshalPKCS1PrivateKey(vk)
		pemtype = PEM_RSA_PRIVATE_KEY

	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(vk)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_EC_PRIVATE_KEY

	case *ed25519.PrivateKey:
		b, err := x509.MarshalPKCS8PrivateKey(vk)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_PRIVATE_KEY

	case *dsa.PrivateKey:
		k := dsaOpenssl{
			Version: 0,
			P:       vk.P,
			Q:       vk.Q,
			G:       vk.G,
			Pub:     vk.Y,
			Priv:    vk.X,
		}
		b, err := asn1.Marshal(k)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_DSA_PRIVATE_KEY
	default:
		return nil, fmt.Errorf("unsupported private key type")
	}
	return &pem.Block{
		Type:  pemtype,
		Bytes: by,
	}, nil
}

type dsaOpenssl struct {
	Version int
	P       *big.Int
	Q       *big.Int
	G       *big.Int
	Pub     *big.Int
	Priv    *big.Int
}
