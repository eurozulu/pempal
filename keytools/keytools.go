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
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
)

// PublicKeySha1Hash returns a hex encoded SHA1 hash of the public key
func PublicKeySha1Hash(key crypto.PublicKey) string {
	if key == nil {
		return ""
	}
	pm, err := MarshalPublicKey(key)
	if err != nil {
		log.Println(err)
		return ""
	}
	hash := sha1.New()
	_, _ = hash.Write(pm.Bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

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
// supports rsa, ecdsa, ed25519 and dsa keytracker
func PublicKeyFromPrivate(pk crypto.PrivateKey) crypto.PublicKey {
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

func MarshalPrivateKey(key crypto.PrivateKey) (*pem.Block, error) {
	var pemtype string
	var by []byte

	switch vk := key.(type) {
	case *rsa.PrivateKey:
		by = x509.MarshalPKCS1PrivateKey(vk)
		pemtype = PEM_RSA_PRIVATE_KEY

	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(vk)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_EC_PRIVATE_KEY

	case *ed25519.PrivateKey:
		b, err := x509.MarshalPKCS8PrivateKey(vk)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_PRIVATE_KEY

	case *dsa.PrivateKey:
		k := dsaOpenssl{
			Version: 0,
			P:       vk.P,
			Q:       vk.Q,
			G:       vk.G,
			Pub:     vk.Y,
			Priv:    vk.X,
		}
		b, err := asn1.Marshal(k)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_PRIVATE_KEY
	default:
		return nil, fmt.Errorf("unsupported private key type")
	}
	return &pem.Block{
		Type:  pemtype,
		Bytes: by,
	}, nil
}

func MarshalPublicKey(key crypto.PublicKey) (*pem.Block, error) {
	var pemtype string
	var by []byte

	switch vk := key.(type) {
	case *rsa.PublicKey:
		by = x509.MarshalPKCS1PublicKey(vk)
		pemtype = PEM_RSA_PUBLIC_KEY

	case *ecdsa.PublicKey:
		b, err := x509.MarshalPKIXPublicKey(vk)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_EC_PUBLIC_KEY
	case *ed25519.PublicKey:
		b, err := x509.MarshalPKIXPublicKey(vk)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_PUBLIC_KEY

	case *dsa.PrivateKey:
		// TODO: Check this is right.
		k := dsaOpenssl{
			Version: 0,
			P:       vk.P,
			Q:       vk.Q,
			G:       vk.G,
			Pub:     vk.Y,
		}
		b, err := asn1.Marshal(k)
		if err != nil {
			return nil, err
		}
		by = b
		pemtype = PEM_PUBLIC_KEY
	default:
		return nil, fmt.Errorf("unsupported private key type")
	}
	return &pem.Block{
		Type:  pemtype,
		Bytes: by,
	}, nil
}

func ParsePrivateKey(blk *pem.Block) (crypto.PrivateKey, error) {
	var prk crypto.PrivateKey
	var err error
	switch blk.Type {

	case PEM_RSA_PRIVATE_KEY:
		prk, err = x509.ParsePKCS1PrivateKey(blk.Bytes)

	case PEM_EC_PRIVATE_KEY:
		prk, err = x509.ParseECPrivateKey(blk.Bytes)

	default:
		//PEM_ANY_PRIVATE_KEY,
		//PEM_DSA_PRIVATE_KEY,
		//PEM_ENCRYPTED_PRIVATE_KEY,
		//PEM_PRIVATE_KEY,
		prk, err = x509.ParsePKCS8PrivateKey(blk.Bytes)
	}
	return prk, err
}

func ParsePublicKey(blk *pem.Block) (crypto.PublicKey, error) {
	var puk crypto.PublicKey
	var err error
	switch blk.Type {

	case PEM_RSA_PUBLIC_KEY:
		puk, err = x509.ParsePKCS1PublicKey(blk.Bytes)

	default:
		puk, err = x509.ParsePKIXPublicKey(blk.Bytes)
	}
	return puk, err
}

type dsaOpenssl struct {
	Version int
	P       *big.Int
	Q       *big.Int
	G       *big.Int
	Pub     *big.Int
	Priv    *big.Int
}
