package identity

import (
	"crypto/x509"
)

type User interface {
	Certificate() *x509.Certificate
	Key() Key
}

type user struct {
	cert *x509.Certificate
	key  Key
}

func (u user) Certificate() *x509.Certificate {
	return u.cert
}

func (u user) Key() Key {
	return u.key
}
