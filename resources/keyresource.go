package resources

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type keyResource struct {
	block *pem.Block
}

func (k keyResource) ResourceType() string {
	return k.block.Type
}

func (k keyResource) MarshalBinary() (data []byte, err error) {
	if k.block == nil {
		return nil, nil
	}
	return k.block.Bytes, nil
}

func (k *keyResource) UnmarshalBinary(data []byte) error {
	key, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return err
	}
	t := "PRIVATE KEY"
	pka := publicKeyAlgorithm(key)
	if pka != x509.UnknownPublicKeyAlgorithm {
		t = fmt.Sprintf("PRIVATE %s KEY", pka.String())
	}
	k.block = &pem.Block{
		Type:  t,
		Bytes: data,
	}
	return nil
}

func (k keyResource) MarshalPEM() (data []byte, err error) {
	if k.block == nil {
		return nil, err
	}
	return pem.EncodeToMemory(k.block), nil
}

func (k keyResource) UnmarshalPem(data []byte) error {
	blk, data := readNextPemBlock("PRIVATE", data)
	if blk == nil {
		return fmt.Errorf("no CERTIFICATE found")
	}
	if err := k.validateKey(blk.Bytes); err != nil {
		return err
	}
	k.block = blk
	return nil
}

func (k keyResource) validateKey(data []byte) error {
	_, err := x509.ParsePKCS8PrivateKey(data)
	return err
}

func publicKeyAlgorithm(prk crypto.PrivateKey) x509.PublicKeyAlgorithm {
	switch _ := prk.(type) {
	case *rsa.PrivateKey:
		return x509.RSA
	case *ecdsa.PrivateKey:
		return x509.ECDSA
	case *ed25519.PrivateKey:
		return x509.Ed25519
	default:
		return x509.UnknownPublicKeyAlgorithm
	}
}
