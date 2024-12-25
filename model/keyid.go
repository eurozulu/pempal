package model

import (
	"bytes"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
)

type KeyId []byte

func NewKeyIdFromKey(key crypto.PublicKey) (KeyId, error) {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	return NewKeyIdFromDer(der)
}

func NewKeyIdFromDer(der []byte) (KeyId, error) {
	hash := sha1.Sum(der)
	return hash[:], nil
}

func (k KeyId) Equals(oid KeyId) bool {
	return bytes.Equal(k, oid)
}

func (k KeyId) String() string {
	return hex.EncodeToString(k)
}

func (k KeyId) MarshalText() (text []byte, err error) {
	return []byte(k.String()), nil
}

func (k *KeyId) UnmarshalText(text []byte) error {
	data, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	*k = data
	return nil
}
