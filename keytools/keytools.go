package keytools

import (
	"bytes"
	"crypto"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"pempal/fileformats"
	"reflect"
	"strconv"
	"strings"
)

// PublicKeySha1Hash returns a hex encoded SHA1 hash of the public key
func PublicKeySha1Hash(key crypto.PublicKey) string {
	if key == nil {
		return ""
	}
	pb, err := fileformats.MarshalPublicKey(key)
	if err != nil {
		log.Println(err)
		return ""
	}
	return SHA1HashString(pb.Bytes)
}

func SHA1HashString(by []byte) string {
	hash := sha1.New()
	_, _ = hash.Write(by)
	return hex.EncodeToString(hash.Sum(nil))
}

// PublicKeyAlgoNames is (unsurprisingly) the names of all the Public Key Algorithms
var PublicKeyAlgoNames = [...]string{
	x509.UnknownPublicKeyAlgorithm: "",
	x509.RSA:                       "RSA",
	x509.DSA:                       "DSA",
	x509.ECDSA:                     "ECDSA",
	x509.Ed25519:                   "Ed25519",
}

// ParsePublicKeyAlgorithm parses the given string into a public key algo
func ParsePublicKeyAlgorithm(s string) x509.PublicKeyAlgorithm {
	for i, pka := range PublicKeyAlgoNames {
		if strings.EqualFold(pka, s) {
			return x509.PublicKeyAlgorithm(i)
		}
	}
	return x509.UnknownPublicKeyAlgorithm
}

// PublicKeyAlgorithm gets the PKA of the given public key
func PublicKeyAlgorithm(pk crypto.PublicKey) x509.PublicKeyAlgorithm {
	if pk == nil {
		return x509.UnknownPublicKeyAlgorithm
	}
	switch pk.(type) {
	case *rsa.PublicKey:
		return x509.RSA
	case *ecdsa.PublicKey:
		return x509.ECDSA
	case *ed25519.PublicKey:
		return x509.Ed25519
	case *dsa.PublicKey:
		return x509.DSA
	default:
		return x509.UnknownPublicKeyAlgorithm
	}
}

// PublicKeyLength returns a representation of the key length.
// Depending on the key typem this can be bitlength (rsa), curve size (ecdsa).  ed25519 has no such length so returns empty string.
func PublicKeyLength(pk crypto.PublicKey) string {
	switch v := pk.(type) {
	case *rsa.PublicKey:
		return strconv.Itoa(v.N.BitLen())
	case *ecdsa.PublicKey:
		return strconv.Itoa(v.Curve.Params().BitSize)
	case *ed25519.PublicKey:
		return ""
	case *dsa.PublicKey:
		return fmt.Sprintf("%d", v.Y)
	default:
		return ""
	}
}

// PublicKeyFromPrivate returns the public key element of the given private key
// supports rsa, ecdsa, ed25519 and dsa keys
func PublicKeyFromPrivate(pk crypto.PrivateKey) crypto.PublicKey {
	if pk == nil || reflect.ValueOf(pk).IsNil() {
		return nil
	}
	switch v := pk.(type) {
	case *rsa.PrivateKey:
		return v.Public()
	case *ecdsa.PrivateKey:
		return v.Public()
	case *ed25519.PrivateKey:
		return v.Public()
	case *dsa.PrivateKey:
		return v.PublicKey
	default:
		return nil
	}
}

// ComparePublicKeys compares two given key for equality
func ComparePublicKeys(pk1 crypto.PublicKey, pk2 crypto.PublicKey) bool {
	switch v := pk1.(type) {
	case *rsa.PublicKey:
		return v.Equal(pk2)
	case *ecdsa.PublicKey:
		return v.Equal(pk2)
	case *ed25519.PublicKey:
		return v.Equal(pk2)
	case *dsa.PublicKey:
		b1, err := x509.MarshalPKIXPublicKey(pk1)
		if err != nil {
			return false
		}
		b2, err := x509.MarshalPKIXPublicKey(pk2)
		if err != nil {
			return false
		}
		return bytes.Equal(b1, b2)

	default:
		return false
	}
}
