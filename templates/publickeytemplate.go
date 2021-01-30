package templates

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/eurozulu/pempal"
	"strings"
)

type PublicKeyTemplate struct {
	PublicKeyAlgorithm   PublicKeyAlgorithm `yaml:"PublicKeyAlgorithm,omitempty"`
	PublicKeyFingerprint string             `yaml:"PublicKeyFingerprint,omitempty"`
	KeyLength            string             `yaml:"KeyLength,omitempty"`
	FilePath             string             `yaml:"Location,omitempty"`
	key                  crypto.PublicKey
}

func NewPublicKeyTemplate(key crypto.PublicKey) *PublicKeyTemplate {
	kt := &PublicKeyTemplate{key: key}
	kt.syncToKey()
	return kt
}

func (t PublicKeyTemplate) Key() crypto.PublicKey {
	return t.key
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
	t.syncToKey()
	return nil
}

func (t *PublicKeyTemplate) MarshalBinary() (data []byte, err error) {
	if t.key == nil {
		return nil, fmt.Errorf("key not set")
	}
	return x509.MarshalPKIXPublicKey(t.key)
}

func (t *PublicKeyTemplate) syncToKey() error {
	if t.key == nil {
		return fmt.Errorf("public key is not set on template")
	}
	t.PublicKeyAlgorithm = PublicKeyAlgorithm(pempal.PublicKeyAlgorithm(t.key))
	t.KeyLength = pempal.PublicKeyLength(t.key)

	by, err := x509.MarshalPKIXPublicKey(t.key)
	if err != nil {
		return err
	}
	t.PublicKeyFingerprint = fingerprint(by)
	return nil
}
