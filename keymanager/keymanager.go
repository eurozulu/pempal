package keymanager

import (
	"crypto"
	"crypto/x509"
)

type KeyManager interface {
	PublicKeys() []Identity
	PrivateKeys() map[Identity]crypto.PrivateKey

	PrivateKey(id Identity) (crypto.PrivateKey, error)
	CertificatesForIdentity(id Identity) []*x509.Certificate
}
