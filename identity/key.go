package identity

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/utils"
	"strings"
)

const pemCipher = x509.PEMCipherAES256

type Key interface {
	Identity() Identity
	String() string
	PublicKey() crypto.PublicKey
	PrivateKey() crypto.PrivateKey
	IsEncrypted() bool
	Encrypt(password []byte) (Key, error)
	Decrypt(password []byte) (Key, error)
	Location() string
}

type key struct {
	prk      *pem.Block
	puk      *pem.Block
	location string
}

func (k key) Identity() Identity {
	return Identity(k.String())
}

func (k key) String() string {
	if k.prk == nil {
		return ""
	}
	s := string(pem.EncodeToMemory(k.prk))
	if k.puk != nil {
		s = strings.Join([]string{s, string(pem.EncodeToMemory(k.puk))}, "")
	}
	return s
}

func (k key) Location() string {
	return k.location
}

func (k key) PrivateKey() crypto.PrivateKey {
	if k.prk == nil {
		return nil
	}
	prk, err := x509.ParsePKCS8PrivateKey(k.prk.Bytes)
	if err != nil {
		logger.Error("Failed to parse private key %v", err)
		return nil
	}
	return prk
}

func (k key) PublicKey() crypto.PublicKey {
	if k.puk == nil {
		return nil
	}
	puk, err := x509.ParsePKIXPublicKey(k.puk.Bytes)
	if err != nil {
		logger.Error("Failed to parse public key %v", err)
		return nil
	}
	return puk
}

func (k key) IsEncrypted() bool {
	if k.prk == nil {
		return false
	}
	return x509.IsEncryptedPEMBlock(k.prk)
}

func (k key) Decrypt(password []byte) (Key, error) {
	if !k.IsEncrypted() {
		return k, nil
	}
	der, err := x509.DecryptPEMBlock(k.prk, password)
	if err != nil {
		return nil, err
	}
	return NewKey(k.location, pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	}))
}

func (k key) Encrypt(password []byte) (Key, error) {
	if k.prk == nil {
		return nil, fmt.Errorf("key is empty and can not be encrypted")
	}
	if k.IsEncrypted() {
		return nil, fmt.Errorf("key already encrypted")
	}
	blk, err := x509.EncryptPEMBlock(rand.Reader, k.prk.Type, k.prk.Bytes, password, pemCipher)
	if err != nil {
		return nil, err
	}
	return &key{
		location: k.location,
		prk:      blk,
		puk:      k.puk,
	}, nil
}

func NewKey(location string, pemdata []byte) (Key, error) {
	k := &key{location: location}
	for len(pemdata) > 0 {
		blk, rest := pem.Decode(pemdata)
		if blk == nil {
			break
		}
		switch resources.ParsePEMType(blk.Type) {
		case resources.PrivateKey:
			if k.prk != nil {
				return nil, fmt.Errorf("multiple private keys found in same pem")
			}
			k.prk = blk
		case resources.PublicKey:
			if k.puk != nil {
				return nil, fmt.Errorf("multiple public keys found in same pem")
			}
			k.puk = blk

		default:
			// ignore non key pems
		}
		pemdata = rest
	}
	if k.prk == nil {
		return nil, fmt.Errorf("no private key found in pem")
	}
	if k.puk == nil {
		// attempt to look up the public key
		if puk, err := utils.PublicKeyFromPrivate(k.PrivateKey()); err != nil {
			logger.Warning("failed to read keys public key %s", err)
		} else {
			k.puk, err = utils.PublicKeyToPEM(puk)
			if err != nil {
				return nil, err
			}
		}
	}
	return k, nil
}
