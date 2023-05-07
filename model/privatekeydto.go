package model

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"pempal/utils"
)

const defaultEncryptCipher = x509.PEMCipherAES256

type PrivateKeyDTO struct {
	PublicKeyAlgorithm string `yaml:"public-key-algorithm"`
	PrivateKey         string `yaml:"private-key,omitempty"`
	PublicKey          string `yaml:"public-key,omitempty"`
	IsEncrypted        bool   `yaml:"is-encrypted,omitempty"`
	KeyParam           string `yaml:"key-param"`

	Identity     string `yaml:"identity"`
	ResourceType string `yaml:"resource-type"`
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
	blk, err := pkr.ToPEMBlock()
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS8PrivateKey(blk.Bytes)
}

func (pkr PrivateKeyDTO) ToPublicKey() (crypto.PublicKey, error) {
	blk, _ := pem.Decode([]byte(pkr.PublicKey))
	if blk == nil {
		return nil, fmt.Errorf("failed to decode public key")
	}
	return x509.ParsePKIXPublicKey(blk.Bytes)
}

func (pkr PrivateKeyDTO) ToPEMBlock() (*pem.Block, error) {
	blk, _ := pem.Decode([]byte(pkr.PrivateKey))
	//pem.Decode([]byte(strings.Trim(pkr.PrivateKey, "`")))
	if blk == nil {
		return nil, fmt.Errorf("failed to decode public key")
	}
	return blk, nil
}

func (pkr *PrivateKeyDTO) Encrypt(pwd []byte) error {
	blk, _ := pem.Decode([]byte(pkr.PrivateKey))
	if blk == nil {
		return fmt.Errorf("failed to decode public key")
	}
	epem, err := x509.EncryptPEMBlock(rand.Reader, blk.Type, blk.Bytes, pwd, defaultEncryptCipher)
	if err != nil {
		return err
	}
	pkr.PrivateKey = string(pem.EncodeToMemory(epem))
	pkr.IsEncrypted = true
	return nil
}

func (pkr *PrivateKeyDTO) Decrypt(pwd []byte) error {
	blk, _ := pem.Decode([]byte(pkr.PrivateKey))
	if blk == nil {
		return fmt.Errorf("failed to decode public key")
	}
	der, err := x509.DecryptPEMBlock(blk, pwd)
	if err != nil {
		return err
	}
	return pkr.UnmarshalBinary(der)
}

func (pkr *PrivateKeyDTO) UnmarshalBinary(der []byte) error {
	pkr.ResourceType = PrivateKey.String()
	pkr.PublicKey = ""
	pkr.IsEncrypted = false
	pkr.PrivateKey = keyPem(PrivateKey, der)

	prk, err := pkr.ToPrivateKey()
	if err != nil {
		return err
	}
	puk, err := utils.PublicKeyFromPrivate(prk)
	if err != nil {
		return err
	}
	if err = pkr.setPublicKey(puk); err != nil {
		return err
	}
	return nil
}

func (pkr *PrivateKeyDTO) UnmarshalPEM(data []byte) error {
	pems := readPEMBlocks(data, PrivateKey)
	if len(pems) != 1 {
		if len(pems) == 0 {
			return fmt.Errorf("failed to find pem private key")
		}
		return fmt.Errorf("multiple private keys found")
	}
	blk := pems[0]
	if x509.IsEncryptedPEMBlock(blk) {
		return pkr.unmarshalEncryptedPEM(data)
	}

	if ParsePEMType(blk.Type) != PrivateKey {
		return fmt.Errorf("Not a %s", PrivateKey.String())
	}
	return pkr.UnmarshalBinary(blk.Bytes)
}

func (pkr *PrivateKeyDTO) unmarshalEncryptedPEM(data []byte) error {
	pkr.ResourceType = PrivateKey.String()
	pkr.PublicKey = ""
	pkr.Identity = ""
	pkr.IsEncrypted = true

	prkBlk := readPEMBlocks(data, PrivateKey)[0] // already checked size in caller
	pkr.PrivateKey = string(pem.EncodeToMemory(prkBlk))

	// Check if the same data contains a single public key
	pukBlks := readPEMBlocks(data, PublicKey)
	if len(pukBlks) == 1 {
		puk, err := x509.ParsePKCS8PrivateKey(pukBlks[0].Bytes)
		if err != nil {
			return err
		}
		return pkr.setPublicKey(puk)
	}
	return nil
}

func (pkr *PrivateKeyDTO) setPublicKey(puk crypto.PublicKey) error {
	pkr.PublicKeyAlgorithm = utils.PublicKeyAlgorithmFromKey(puk).String()
	pukder, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return err
	}
	pkr.PublicKey = keyPem(PublicKey, pukder)
	pkr.Identity = pkr.String()

	return nil
}

func keyPem(rt ResourceType, der []byte) string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  rt.PEMString(),
		Bytes: der,
	}))
}
