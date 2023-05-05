package keymanager

import (
	"crypto"
	"crypto/x509"
)

type User interface {
	Identity() Identity
	Certificate() *x509.Certificate
	Key() crypto.PrivateKey
}

type user struct {
	id   Identity
	cert *x509.Certificate
	key  crypto.PrivateKey
}

func (u user) Identity() Identity {
	return u.id
}

func (u user) Certificate() *x509.Certificate {
	return u.cert
}

func (u user) Key() crypto.PrivateKey {
	return u.key
}
