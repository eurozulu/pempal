package fileformats

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"pempal/keytools"
)

type derKeyFormat struct{}

func (d derKeyFormat) Format(by []byte) ([]*pem.Block, error) {
	var puk crypto.PublicKey

	// First attempt to parse as a private key
	prk, err := x509.ParsePKCS8PrivateKey(by)
	if err != nil {
		// not a private key, try as a public
		puk, _ = x509.ParsePKIXPublicKey(by)

	} else {
		// is a private key, derive puk from it
		puk = keytools.PublicKeyFromPrivate(prk)
	}
	if puk == nil {
		// Not a known der format, try as PEM
		return pemFormatter.Format(by)
	}

	ka := keytools.PublicKeyAlgorithm(puk)
	if ka == x509.UnknownPublicKeyAlgorithm {
		return nil, fmt.Errorf(ka.String())
	}

	pb, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return nil, err
	}

	blocks := []*pem.Block{{
		Type:  PublicKeyPEMType(ka),
		Bytes: pb,
	}}

	if prk != nil {
		blocks = append(blocks, &pem.Block{
			Type:  PrivateKeyPEMType(ka),
			Bytes: by,
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
	default:
		return "PRIVATE KEY"
	}
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
