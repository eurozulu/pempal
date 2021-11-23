package keytools

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
)

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

func ParsePrivateKey(blk *pem.Block) (crypto.PrivateKey, error) {
	var prk crypto.PrivateKey
	var err error
	switch blk.Type {

	case PEM_RSA_PRIVATE_KEY:
		prk, err = x509.ParsePKCS1PrivateKey(blk.Bytes)

	case PEM_EC_PRIVATE_KEY:
		prk, err = x509.ParseECPrivateKey(blk.Bytes)

	default:
		//PEM_ANY_PRIVATE_KEY,
		//PEM_DSA_PRIVATE_KEY,
		//PEM_ENCRYPTED_PRIVATE_KEY,
		//PEM_PRIVATE_KEY,
		prk, err = x509.ParsePKCS8PrivateKey(blk.Bytes)
	}
	return prk, err
}

func ParsePublicKey(blk *pem.Block) (crypto.PublicKey, error) {
	var puk crypto.PublicKey
	var err error
	switch blk.Type {

	case PEM_RSA_PUBLIC_KEY:
		puk, err = x509.ParsePKCS1PublicKey(blk.Bytes)

	default:
		puk, err = x509.ParsePKIXPublicKey(blk.Bytes)
	}
	return puk, err
}

type dsaOpenssl struct {
	Version int
	P       *big.Int
	Q       *big.Int
	G       *big.Int
	Pub     *big.Int
	Priv    *big.Int
}
