package model

import "crypto"

type Issuer struct {
	cert *Certificate
	key  *PrivateKey
}

func (i Issuer) String() string {
	return i.cert.Subject.String()
}

func (i *Issuer) PublicKey() *PublicKey {
	return i.key.Public()
}

func (i *Issuer) PrivateKey() *PrivateKey {
	return i.key
}

func (i *Issuer) Signer() crypto.Signer {
	return i.key.Signer()
}

func (i Issuer) Certificate() *Certificate {
	return i.cert
}

func NewIssuer(c *Certificate, k *PrivateKey) *Issuer {
	return &Issuer{
		cert: c,
		key:  k,
	}
}
