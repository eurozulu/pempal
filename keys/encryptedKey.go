package keys

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"io"
	"strings"
)

// encryptedKey represents an encrypted private key
type encryptedKey struct {
	location string
	pemBlock *pem.Block
}

func (e encryptedKey) PublicKey() crypto.PublicKey {
	return nil
}

func (e encryptedKey) privateKey() crypto.PrivateKey {
	return nil
}

func (e encryptedKey) PublicKeyAlgorithm() x509.PublicKeyAlgorithm {
	// attempt to find the key type from the PEM description
	blkType := e.pemBlock.Type
	if strings.Contains(blkType, x509.RSA.String()) {
		return x509.RSA
	}
	if strings.Contains(blkType, x509.ECDSA.String()) {
		return x509.ECDSA
	}
	if strings.Contains(blkType, x509.Ed25519.String()) {
		return x509.Ed25519
	}
	return x509.UnknownPublicKeyAlgorithm
}

func (e encryptedKey) Location() string {
	return e.location
}

// DecryptKey decrypts this key into an unencrypted derKey
func (e encryptedKey) DecryptKey(password []byte) (Key, error) {
	der, err := x509.DecryptPEMBlock(e.pemBlock, password)
	if err != nil {
		return nil, err
	}
	return parseDERKey(e.location, der)
}

func (e encryptedKey) WriteKey(out io.Writer) error {
	return pem.Encode(out, e.pemBlock)
}

func IsEncrypted(k Key) bool {
	_, ok := k.(*encryptedKey)
	return ok
}
