package model

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

type PublicKey struct {
	puk crypto.PublicKey
}

func (k PublicKey) String() string {
	return fmt.Sprintf("%s\t%s", k.Fingerprint().String(), k.PublicKeyAlgorithm().String())
}

func (k PublicKey) Fingerprint() Fingerprint {
	data, err := k.MarshalBinary()
	if err != nil {
		return Fingerprint{}
	}
	return NewFingerPrint(data)
}

func (k PublicKey) ResourceType() ResourceType {
	return ResourceTypePublicKey
}

func (k PublicKey) Public() crypto.PublicKey {
	return k.puk
}

func (k PublicKey) PublicKeyAlgorithm() PublicKeyAlgorithm {
	switch k.puk.(type) {
	case *rsa.PublicKey:
		return PublicKeyAlgorithm(x509.RSA)
	case *ecdsa.PublicKey:
		return PublicKeyAlgorithm(x509.ECDSA)
	case *ed25519.PublicKey:
		return PublicKeyAlgorithm(x509.Ed25519)
	case *dsa.PublicKey:
		return PublicKeyAlgorithm(x509.DSA)
	default:
		return PublicKeyAlgorithm(x509.UnknownPublicKeyAlgorithm)
	}
}

func (k *PublicKey) MarshalBinary() (data []byte, err error) {
	return x509.MarshalPKIXPublicKey(k.puk)
}

func (k *PublicKey) UnmarshalBinary(data []byte) error {
	puk, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return err
	}
	k.puk = puk
	return nil
}

func (k *PublicKey) MarshalText() (text []byte, err error) {
	der, err := k.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: ResourceTypePublicKey.String(), Bytes: der}), nil
}

func (k *PublicKey) UnmarshalText(text []byte) error {
	var der []byte
	for len(text) > 0 {
		blk, rest := pem.Decode(text)
		if blk == nil {
			break
		}
		if strings.Contains(blk.Type, ResourceTypePublicKey.String()) {
			der = blk.Bytes
			break
		}
		text = rest
	}
	if der == nil {
		return errors.New("no public key PEM found")
	}
	return k.UnmarshalBinary(der)
}

func NewPublicKeyFromPem(blk *pem.Block) (*PublicKey, error) {
	puk := &PublicKey{}
	if err := puk.UnmarshalText(pem.EncodeToMemory(blk)); err != nil {
		return nil, err
	}
	return puk, nil
}

func NewPublicKey(puk crypto.PublicKey) *PublicKey {
	return &PublicKey{puk: puk}
}
