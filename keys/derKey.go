package keys

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
)

// used to encrypt the key
const pemCipher = x509.PEMCipherAES256

// derKey represents an unencrypted private key
type derKey struct {
	location string
	pk       crypto.PrivateKey
}

func (d derKey) PublicKey() crypto.PublicKey {
	switch v := d.pk.(type) {
	case *rsa.PrivateKey:
		return v.PublicKey
	case *ecdsa.PrivateKey:
		return v.PublicKey
	case *ed25519.PrivateKey:
		return v.Public()
	default:
		return nil
	}
}

func (d derKey) privateKey() crypto.PrivateKey {
	x509.CreateCertificateRequest()
	return d.pk
}

func (d derKey) PublicKeyAlgorithm() x509.PublicKeyAlgorithm {
	switch _ := d.pk.(type) {
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

func (d derKey) Location() string {
	return d.location
}

func (d derKey) WriteKey(out io.Writer) error {
	buf := bytes.NewBuffer(nil)
	if err := d.WriteKeyDER(buf); err != nil {
		return err
	}
	return pem.Encode(out, &pem.Block{
		Type:  fmt.Sprintf("PRIVATE %s KEY", d.PublicKeyAlgorithm().String()),
		Bytes: buf.Bytes(),
	})
}

func (d derKey) WriteKeyDER(out io.Writer) error {
	der, err := x509.MarshalPKCS8PrivateKey(d.PublicKey())
	if err != nil {
		return err
	}
	_, err = out.Write(der)
	return err
}

// EncryptKey encrypts this key into an encryptedKey using the given password
func (d derKey) EncryptKey(password []byte) (Key, error) {
	der, err := x509.MarshalPKCS8PrivateKey(d.pk)
	if err != nil {
		return nil, err
	}
	blk, err := x509.EncryptPEMBlock(rand.Reader,
		fmt.Sprintf("PRIVATE %s KEY", d.PublicKeyAlgorithm().String()),
		der,
		password,
		pemCipher)
	if err != nil {
		return nil, err
	}
	return &encryptedKey{
		location: d.location,
		pemBlock: blk,
	}, nil
}
