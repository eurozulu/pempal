package keys

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/utils"
)

type Key interface {
	Identity() model.Identity
	PublicKey() crypto.PublicKey
	PrivateKey() crypto.PrivateKey
	IsEncrypted() bool
}

type key struct {
	id  model.Identity
	prk *pem.Block
}

func (k key) Identity() model.Identity {
	return k.id
}

func (k key) PublicKey() crypto.PublicKey {
	puk, err := k.Identity().PublicKey()
	if err != nil {
		logger.Error("Failed to parse public key %v", err)
		return nil
	}
	return puk
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

func (k key) IsEncrypted() bool {
	return x509.IsEncryptedPEMBlock(k.prk)
}

func NewKeyPair(keypem *pem.Block, puk crypto.PublicKey) (Key, error) {
	id, err := model.NewIdentity(puk)
	if err != nil {
		return nil, err
	}
	return &key{
		id:  id,
		prk: keypem,
	}, nil
}

func NewKey(keypem *pem.Block) (Key, error) {
	k := &key{
		prk: keypem,
	}
	prk := k.PrivateKey()
	if prk == nil {
		return nil, fmt.Errorf("failed to parse key")
	}
	puk, err := utils.PublicKeyFromPrivate(prk)
	if err != nil {
		return nil, err
	}
	return NewKeyPair(keypem, puk)
}
