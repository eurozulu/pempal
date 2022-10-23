package keys

import (
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

var hasher = sha256.New()

// Key represents a private key
type Key interface {
	PublicKey() crypto.PublicKey
	PublicKeyAlgorithm() x509.PublicKeyAlgorithm
	Location() string
	WriteKey(out io.Writer) error
	privateKey() crypto.PrivateKey
}

// parseDERKey attempts to parse the given der byte block into a private key
func parseDERKey(location string, der []byte) (Key, error) {
	pk, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}
	return &derKey{
		location: location,
		pk:       pk,
	}, nil
}

// parsePEMKey attempts to parse the given pem byte block into a private key
// If there are more than one PEM blocks in the bytes, locates the first private key.
func parsePEMKey(location string, by []byte) (Key, error) {
	var blk *pem.Block
	for {
		blk, by = pem.Decode(by)
		if blk == nil {
			return nil, fmt.Errorf("no key pem found")
		}
		if !strings.Contains(blk.Type, "KEY") {
			// skip non key blocks in the same file
			continue
		}

		if x509.IsEncryptedPEMBlock(blk) {
			return &encryptedKey{
				location: location,
				pemBlock: blk,
			}, nil
		}
		// parse as an unencrypted der key
		return parseDERKey(location, blk.Bytes)
	}
}

// isPem checks if the given byte block is a PEM foramt
func isPem(by []byte) bool {
	s := string(by)
	return strings.Contains(s, "-----BEGIN") && strings.Contains(s, "-----END")
}

// NewKey attemps to load a private key from the given file location
func NewKey(location string) (Key, error) {
	by, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	if isPem(by) {
		return parsePEMKey(location, by)
	}
	// parse as an unencrypted der key
	return parseDERKey(location, by)
}

// PublicKeyHash generates the hash of the Public key for the given private key
func PublicKeyHash(puk crypto.PublicKey) ([]byte, error) {
	if puk == nil {
		return nil, fmt.Errorf("no public key available.")
	}
	by, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return nil, err
	}
	hasher.Reset()
	hasher.Write(by)
	return hasher.Sum(nil), nil
}
