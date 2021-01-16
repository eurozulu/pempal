package templates

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal"
	"golang.org/x/crypto/ssh"
)

type PublicKeyTemplate struct {
	PublicKeyAlgorithm   PublicKeyAlgorithm `yaml:"PublicKeyAlgorithm,omitempty"`
	PublicKeyFingerprint string             `yaml:"PublicKeyFingerprint,omitempty"`
	KeyLength            string             `yaml:"KeyLength,omitempty"`
	FilePath             string             `yaml:"-"`
	key                  crypto.PublicKey
}

func (t *PublicKeyTemplate) String() string {
	return TemplateString(t)
}

func (t *PublicKeyTemplate) Location() string {
	return t.FilePath
}

func (t *PublicKeyTemplate) UnmarshalBinary(by []byte) error {
	k, err := x509.ParsePKIXPublicKey(by)
	if err != nil {
		return err
	}
	t.key = k
	t.PublicKeyAlgorithm = PublicKeyAlgorithm(pempal.PublicKeyAlgorithm(k))
	sk, err := ssh.NewPublicKey(k)
	if err != nil {
		return err
	}
	t.PublicKeyFingerprint = ssh.FingerprintSHA256(sk)
	t.KeyLength = pempal.PublicKeyLength(k)

	return nil
}

func (t *PublicKeyTemplate) MarshalBinary() (data []byte, err error) {
	if t.key == nil {
		return nil, fmt.Errorf("key not set")
	}
	return x509.MarshalPKIXPublicKey(t.key)
}

func (t *PublicKeyTemplate) UnmarshalPEM(bl *pem.Block) error {
	return t.UnmarshalBinary(bl.Bytes)
}

func (t *PublicKeyTemplate) MarshalPEM() (*pem.Block, error) {
	by, err := t.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: by,
	}, nil
}

type SSHPublicKeyTemplate struct {
	KeyType              string `yaml:"KeyType,omitempty"`
	PublicKeyFingerprint string `yaml:"PublicKeyFingerprint,omitempty"`
	Comment              string `yaml:"Comment,omitempty"`
	FilePath             string `yaml:"-"`

	key ssh.PublicKey
}

func (t *SSHPublicKeyTemplate) String() string {
	return TemplateString(t)
}

func (t SSHPublicKeyTemplate) Location() string {
	return t.FilePath
}

func (t *SSHPublicKeyTemplate) UnmarshalBinary(data []byte) error {
	k, c, _, _, err := ssh.ParseAuthorizedKey(data)
	if err != nil {
		return err
	}
	t.key = k
	t.KeyType = k.Type()
	t.PublicKeyFingerprint = ssh.FingerprintSHA256(k)
	t.Comment = c
	return nil
}

func (t *SSHPublicKeyTemplate) MarshalBinary() (data []byte, err error) {
	if t.key == nil {
		return nil, fmt.Errorf("no binary private key data available")
	}
	return ssh.MarshalAuthorizedKey(t.key), nil
}
