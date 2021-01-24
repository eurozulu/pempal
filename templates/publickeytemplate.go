package templates

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal"
	"golang.org/x/crypto/ssh"
	"strings"
)

type PublicKeyTemplate struct {
	PublicKeyAlgorithm   PublicKeyAlgorithm `yaml:"PublicKeyAlgorithm,omitempty"`
	PublicKeyFingerprint string             `yaml:"PublicKeyFingerprint,omitempty"`
	KeyLength            string             `yaml:"KeyLength,omitempty"`
	FilePath             string             `yaml:"Location,omitempty"`
	key                  crypto.PublicKey
}

func (t *PublicKeyTemplate) String() string {
	pka := t.PublicKeyAlgorithm.String()
	if pka == "" {
		pka = "???"
	}
	pkf := t.PublicKeyFingerprint
	if pkf == "" {
		pkf = "???"
	}
	return strings.Join([]string{TemplateType(t), pka, t.KeyLength, pkf, t.Location()}, "\t")
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

type SSHPublicKeyTemplate struct {
	Comment string `yaml:"Comment,omitempty"`
	PublicKeyTemplate
}

func (t *SSHPublicKeyTemplate) String() string {
	s := strings.Split(t.PublicKeyTemplate.String(), "\t")
	if len(s) > 0 {
		s[0] = TemplateType(t)
		s[len(s)-1] = t.Comment
		s = append(s, t.Location())
	}
	return strings.Join(s, "\t")
}

func (t *SSHPublicKeyTemplate) UnmarshalBinary(data []byte) error {
	c, k, err := pempal.ParseSSHPublicKey(data)
	if err != nil {
		return err
	}
	t.key = k
	t.Comment = c
	by, err := x509.MarshalPKIXPublicKey(k)
	if err != nil {
		return err
	}
	return t.PublicKeyTemplate.UnmarshalBinary(by)
}

func (t *SSHPublicKeyTemplate) MarshalBinary() (data []byte, err error) {
	if t.key == nil {
		return nil, fmt.Errorf("no binary private key data available")
	}
	spk, err := ssh.NewPublicKey(t.key)
	if err != nil {
		return nil, err
	}
	by := ssh.MarshalAuthorizedKey(spk)
	if t.Comment != "" {
		by = append(by, ' ')
		by = append(by, []byte(t.Comment)...)
	}
	return by, nil
}
