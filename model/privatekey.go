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
	"github.com/eurozulu/pempal/logging"
	"strings"
)

type PrivateKey struct {
	prk crypto.PrivateKey
}

func (k PrivateKey) ResourceType() ResourceType {
	return ResourceTypePrivateKey
}

func (k PrivateKey) Fingerprint() Fingerprint {
	data, err := k.MarshalBinary()
	if err != nil {
		logging.Error("Error marshalling key: %v", err)
		return Fingerprint{}
	}
	return NewFingerPrint(data)
}

func (k PrivateKey) String() string {
	var s string
	switch x509.PublicKeyAlgorithm(k.PublicKeyAlgorithm()) {
	case x509.RSA:
		s = fmt.Sprintf("%s (%d)", k.PublicKeyAlgorithm().String(), k.RSAKeyLength())
	case x509.ECDSA:
		s = fmt.Sprintf("%s (%s)", k.PublicKeyAlgorithm().String(), k.ECDSACurve().String())
	default:
		s = k.PublicKeyAlgorithm().String()
	}
	return fmt.Sprintf("%s\t%s", k.Fingerprint(), s)
}

func (k PrivateKey) RSAKeyLength() int {
	key, ok := k.prk.(*rsa.PrivateKey)
	if !ok {
		return 0
	}
	return key.N.BitLen()
}

func (k PrivateKey) ECDSACurve() EllipticCurve {
	key, ok := k.prk.(*ecdsa.PrivateKey)
	if !ok {
		return 0
	}
	c, err := NewCurve(key.Curve)
	if err != nil {
		logging.Error("Error parsing ecdsa curve: %v", err)
		return 0
	}
	return c
}

func (k PrivateKey) PublicKeyAlgorithm() PublicKeyAlgorithm {
	switch k.prk.(type) {
	case *rsa.PrivateKey:
		return PublicKeyAlgorithm(x509.RSA)
	case *ecdsa.PrivateKey:
		return PublicKeyAlgorithm(x509.ECDSA)
	case ed25519.PrivateKey:
		return PublicKeyAlgorithm(x509.Ed25519)
	case *dsa.PrivateKey:
		return PublicKeyAlgorithm(x509.DSA)
	default:
		return PublicKeyAlgorithm(x509.UnknownPublicKeyAlgorithm)
	}
}

func (k PrivateKey) Private() crypto.PrivateKey {
	return k.prk
}

func (k PrivateKey) Signer() crypto.Signer {
	return k.prk.(crypto.Signer)
}

func (k PrivateKey) Public() *PublicKey {
	return NewPublicKey(k.Signer().Public())
}

func (k *PrivateKey) UnmarshalBinary(der []byte) error {
	prk, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return err
	}
	(*k).prk = prk
	return nil
}

func (k PrivateKey) MarshalBinary() (der []byte, err error) {
	return x509.MarshalPKCS8PrivateKey(k.prk)
}

// UnmarshalText will parse the data as pem block(s), to find a private key.
// if data contains multiple pem blocks, the first private key in the data is used.
func (k *PrivateKey) UnmarshalText(data []byte) error {
	var der []byte
	for len(data) > 0 {
		blk, rest := pem.Decode(data)
		if blk == nil {
			break
		}
		if strings.Contains(blk.Type, ResourceTypePrivateKey.String()) {
			der = blk.Bytes
			break
		}
		data = rest
	}
	if der == nil {
		return errors.New("no private key PEM found")
	}
	return k.UnmarshalBinary(der)
}

// MarshalText marshals the key into a PEM
func (k PrivateKey) MarshalText() (text []byte, err error) {
	der, err := k.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: ResourceTypePrivateKey.String(), Bytes: der}), nil
}

func NewPrivateKeyFromPem(blk *pem.Block) (*PrivateKey, error) {
	prk := &PrivateKey{}
	if err := prk.UnmarshalText(pem.EncodeToMemory(blk)); err != nil {
		return nil, err
	}
	return prk, nil
}

func NewPrivateKey(prk crypto.PrivateKey) *PrivateKey {
	return &PrivateKey{prk: prk}
}
