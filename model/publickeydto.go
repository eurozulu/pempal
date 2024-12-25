package model

import (
	"crypto"
	"crypto/x509"
)

type PublicKeyDTO struct {
	PublicKey crypto.PublicKey `json:"_"`
}

func (p PublicKeyDTO) Equals(t crypto.PublicKey) bool {
	type eq interface {
		Equals(t crypto.PublicKey) bool
	}
	puk, ok := t.(eq)
	if !ok {
		return false
	}
	return puk.Equals(p)
}

func (p PublicKeyDTO) MarshalText() (text []byte, err error) {
	if p.PublicKey == nil {
		return nil, nil
	}
	by, err := x509.MarshalPKIXPublicKey(p.PublicKey)
	if err != nil {
		return nil, err
	}
	return []byte(EncodeAsBase64(by)), nil
}

func (p *PublicKeyDTO) UnmarshalText(text []byte) (err error) {
	by, err := DecodeAsBase64(string(text))
	if err != nil {
		return err
	}
	puk, err := x509.ParsePKIXPublicKey(by)
	if err != nil {
		return err
	}
	p.PublicKey = puk
	return nil
}
