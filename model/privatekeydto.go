package model

import (
	"crypto"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"pempal/utils"
)

type PrivateKeyDTO struct {
	PublicKeyAlgorithm string `yaml:"private-key-algorithm"`
	PrivateKey         string `yaml:"private-key,omitempty"`
	PublicKey          string `yaml:"public-key,omitempty"`
	IsEncrypted        bool   `yaml:"is-encrypted,omitempty"`
	Identity           string `yaml:"identity"`
}

func (pkr PrivateKeyDTO) String() string {
	if pkr.PublicKey == "" {
		return ""
	}
	der, err := hex.DecodeString(pkr.PublicKey)
	if err != nil {
		return ""
	}
	return MD5PublicKey(der)
}

func (pkr PrivateKeyDTO) ToPrivateKey() (crypto.PrivateKey, error) {
	der, err := hex.DecodeString(pkr.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("private-key is invalid hex or empty %v", err)
	}
	return x509.ParsePKCS8PrivateKey(der)
}

func (pkr *PrivateKeyDTO) UnmarshalBinary(data []byte) error {
	pkr.PublicKey = ""
	pkr.IsEncrypted = false
	pkr.PrivateKey = hex.EncodeToString(data)
	prk, err := pkr.ToPrivateKey()
	if err != nil {
		return fmt.Errorf("private-key failed to parse  %v", err)
	}
	puk, err := utils.PublicKeyFromPrivate(prk)
	if err != nil {
		return err
	}
	pkr.PublicKeyAlgorithm = utils.PublicKeyAlgorithmFromKey(puk).String()
	pukder, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return err
	}
	pkr.PublicKey = hex.EncodeToString(pukder)
	pkr.Identity = pkr.String()
	return nil
}
