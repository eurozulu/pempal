package templates

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal"
	"strings"
)

// PrivateKeyTemplate is for private keys encoded in pkcs8 or rsa
type PrivateKeyTemplate struct {
	Encrypted            bool               `yaml:"Encrypted,omitempty"`
	PublicKeyAlgorithm   PublicKeyAlgorithm `yaml:"PublicKeyAlgorithm,omitempty"`
	PublicKeyFingerprint string             `yaml:"PublicKeyFingerprint,omitempty"`
	FilePath             string             `yaml:"Location,omitempty"`
	Passphrase           string             `yaml:"-"`

	key      crypto.PrivateKey
	pemBlock *pem.Block // Keeps pem block for when key is encrypted
}

func (t *PrivateKeyTemplate) Key() crypto.PrivateKey {
	return t.key
}

func (t *PrivateKeyTemplate) PublicKey() crypto.PublicKey {
	return pempal.PublicKeyFromPrivate(t.key)
}

func (t *PrivateKeyTemplate) String() string {
	pkf := t.PublicKeyFingerprint
	if pkf == "" {
		pkf = "???"
	}
	pka := t.PublicKeyAlgorithm.String()
	if pka == "" {
		pka = "???"
	}

	e := "encrypted"
	if !t.Encrypted {
		e = "unencrypted"
	}
	return strings.Join([]string{TemplateType(t), pka, e, pkf, t.Location()}, "\t")
}

func (t *PrivateKeyTemplate) Decrypt(passphrase string) error {
	if t.pemBlock == nil {
		return fmt.Errorf("template %s has no binary pem data", t.FilePath)
	}
	t.Passphrase = passphrase
	return t.unmarshalPEM(t.pemBlock)
}

func (t PrivateKeyTemplate) Location() string {
	return t.FilePath
}

func (t *PrivateKeyTemplate) UnmarshalBinary(data []byte) error {
	bl, _ := pem.Decode(data)
	if bl != nil {
		return t.unmarshalPEM(bl)
	}

	k, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		if t.Encrypted {
			return nil
		}
		// Try to parse as an rsa key
		k, err = x509.ParsePKCS1PrivateKey(data)
		if err != nil {
			return err
		}
	}
	t.key = k
	pk := t.PublicKey()
	by, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return err
	}
	t.PublicKeyFingerprint = fingerprint(by)
	t.PublicKeyAlgorithm = PublicKeyAlgorithm(pempal.PublicKeyAlgorithm(pk))
	return nil
}

func (t *PrivateKeyTemplate) MarshalBinary() (data []byte, err error) {
	if t.key == nil {
		if t.Encrypted {
			return nil, fmt.Errorf("template %s is encrypted and cannot be marshalled without passphrase", t.FilePath)
		}
		return nil, fmt.Errorf("template %s has no binary key data", t.FilePath)
	}
	return x509.MarshalPKCS8PrivateKey(t.key)
}

func (t *PrivateKeyTemplate) unmarshalPEM(bl *pem.Block) error {
	t.Encrypted = x509.IsEncryptedPEMBlock(bl)
	if t.Encrypted && t.Passphrase != "" {
		by, err := x509.DecryptPEMBlock(bl, []byte(t.Passphrase))
		if err != nil {
			return err
		}
		bl.Bytes = by
	}
	t.pemBlock = bl
	return t.UnmarshalBinary(bl.Bytes)
}
