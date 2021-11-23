package keytracker

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"pempal/keytools"
	"pempal/pemreader"
	"strings"
)

const encryptedHash = "*"

// A Key represents a Private key, identified by a "KeyHash"
// The hash is a SHA1 hash of the keys Public key unless the private key is encrypted, in which case
// the hash is a SHA1 hash of the encrypted key itself, preceeded by a '*'
type Key interface {
	fmt.Stringer
	// PublicKey retuens the private keys public key, if available.
	// For encrypted keys this is only available when an associated public key file is present.
	PublicKey() crypto.PublicKey

	// Location is the location of the private key
	Location() string

	// PublicLocation is the location of a matched public key, if available.
	PublicLocation() string

	// Type is the pem type of the private key
	Type() string

	// IsEncrypted checks if the key is encrypted.
	// Returns true is the PEM type is ENCRYPTED_PRIVATE_KEY, or any pem header contains the "ENCRYPT" word.
	IsEncrypted() bool

	// PrivateKey returns the parsed private key, if available
	PrivateKey() (crypto.PrivateKey, error)

	// PrivateKeyDecrypted returns the decrypted private using the given password.
	PrivateKeyDecrypted(passwd string) (crypto.PrivateKey, error)
}

// key represents a private key in its PEM form.
// key has an additional public key field, used for encryoted keys where an associated public key file was found
type key struct {
	pemBlock *pem.Block
	puk      *pem.Block
}

func (k key) PublicKey() crypto.PublicKey {
	// attempt to read the PUK from the private key first.
	prk, _ := k.PrivateKey()
	if prk != nil {
		return keytools.PublicKeyFromPrivate(prk)
	}
	// no private key available (encrypted), return paired puk if available
	if k.puk != nil {
		k, err := keytools.ParsePublicKey(k.puk)
		if err != nil {
			log.Println(err)
		}
		return k
	}
	return nil
}

func (k key) Location() string {
	if len(k.pemBlock.Headers) == 0 {
		return ""
	}
	return k.pemBlock.Headers[pemreader.LocationHeaderKey]
}

func (k key) PublicLocation() string {
	if k.puk == nil || len(k.puk.Headers) == 0 {
		return ""
	}
	return k.puk.Headers[pemreader.LocationHeaderKey]
}

func (k key) IsEncrypted() bool {
	if k.pemBlock.Type == keytools.PEM_ENCRYPTED_PRIVATE_KEY {
		return true
	}
	if x509.IsEncryptedPEMBlock(k.pemBlock) {
		return true
	}
	pt, ok := k.pemBlock.Headers["Proc-Type"]
	return ok && strings.HasSuffix(pt, "ENCRYPTED")
}

func (k key) Type() string {
	return k.pemBlock.Type
}

func (k key) PrivateKey() (crypto.PrivateKey, error) {
	if k.IsEncrypted() {
		return nil, fmt.Errorf("key is encrypted, requires password")
	}
	prk, err := keytools.ParsePrivateKey(k.pemBlock)
	if err != nil {
		return nil, err
	}
	return prk, nil
}

func (k key) PrivateKeyDecrypted(passwd string) (crypto.PrivateKey, error) {
	if !k.IsEncrypted() {
		return k.PrivateKey()
	}
	der, err := x509.DecryptPEMBlock(k.pemBlock, []byte(passwd))
	if err != nil {
		return nil, err
	}
	return keytools.ParsePrivateKey(&pem.Block{
		Type:  k.pemBlock.Type,
		Bytes: der,
	})
}

func (k key) String() string {
	puk := k.PublicKey()
	if puk != nil {
		return keytools.PublicKeySha1Hash(puk)
	}
	return strings.Join([]string{encryptedHash, stringHash(k.pemBlock.Bytes)}, "")
}

func NewKey(blk *pem.Block) Key {
	return &key{
		pemBlock: blk,
	}
}
