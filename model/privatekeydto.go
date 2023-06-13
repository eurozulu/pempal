package model

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
)

const defaultEncryptCipher = x509.PEMCipherAES256

type PrivateKeyDTO struct {
	Id                 string `yaml:"identity,omitempty"`
	PublicKeyAlgorithm string `yaml:"public-key-algorithm"`
	PrivateKey         string `yaml:"private-key,omitempty"`
	PublicKey          string `yaml:"public-key,omitempty"`
	IsEncrypted        bool   `yaml:"is-encrypted,omitempty"`
	KeyParam           string `yaml:"key-param,omitempty"`

	ResourceType string `yaml:"resource-type"`
}

func (p PrivateKeyDTO) String() string {
	return p.PrivateKey
}

func (p *PrivateKeyDTO) UnmarshalPEM(data []byte) error {
	p.reset()
	// scan the pems looking for private AND public keys.
	// If private is encypted and only one public key found, public is assumed to be its pair.
	for len(data) > 0 {
		blk, rest := pem.Decode(data)
		if blk == nil {
			break
		}
		rt := ParsePEMType(blk.Type)
		if rt == PrivateKey {
			if err := p.setPrivateKey(blk); err != nil {
				return fmt.Errorf("failed to unmarshal private key  %v", err)
			}
		} else if rt == PublicKey {
			if err := p.setPublicKey(blk); err != nil {
				return fmt.Errorf("failed to unmarshal public key  %v", err)
			}
		}
		data = rest
		continue
	}
	return nil
}

func (p *PrivateKeyDTO) UnmarshalBinary(der []byte) error {
	p.reset()
	return p.setPrivateKey(&pem.Block{
		Type:  PrivateKey.String(),
		Bytes: der,
	})
}

func (p PrivateKeyDTO) ToPrivateKey() (crypto.PrivateKey, error) {
	blk, _ := pem.Decode([]byte(p.PrivateKey))
	if blk == nil {
		return nil, fmt.Errorf("failed to decode public key")
	}
	return x509.ParsePKCS8PrivateKey(blk.Bytes)
}

func (p PrivateKeyDTO) ToPublicKey() (crypto.PublicKey, error) {
	blk, _ := pem.Decode([]byte(p.PublicKey))
	if blk == nil {
		return nil, fmt.Errorf("failed to decode public key")
	}
	return x509.ParsePKIXPublicKey(blk.Bytes)
}

func (p *PrivateKeyDTO) Encrypt(pwd []byte) error {
	blk, _ := pem.Decode([]byte(p.PrivateKey))
	if blk == nil {
		return fmt.Errorf("failed to decode public key")
	}
	epem, err := x509.EncryptPEMBlock(rand.Reader, blk.Type, blk.Bytes, pwd, defaultEncryptCipher)
	if err != nil {
		return err
	}
	p.PrivateKey = string(pem.EncodeToMemory(epem))
	p.IsEncrypted = true
	return nil
}

func (p *PrivateKeyDTO) Decrypt(pwd []byte) error {
	blk, _ := pem.Decode([]byte(p.PrivateKey))
	if blk == nil {
		return fmt.Errorf("failed to decode public key")
	}
	der, err := x509.DecryptPEMBlock(blk, pwd)
	if err != nil {
		return err
	}
	return p.UnmarshalBinary(der)
}

func (p *PrivateKeyDTO) reset() {
	p.ResourceType = PrivateKey.String()
	p.PublicKeyAlgorithm = x509.UnknownPublicKeyAlgorithm.String()
	p.PrivateKey = ""
	p.PublicKey = ""
	p.IsEncrypted = false
	p.KeyParam = ""
}

func (p *PrivateKeyDTO) setPrivateKey(block *pem.Block) error {
	if p.PrivateKey != "" {
		return fmt.Errorf("multiple private keys found in same resource")
	}
	p.PrivateKey = string(pem.EncodeToMemory(block))
	p.IsEncrypted = x509.IsEncryptedPEMBlock(block)

	if p.IsEncrypted {
		return nil
	}
	prk, err := p.ToPrivateKey()
	if err != nil {
		return err
	}
	puk, err := utils.PublicKeyFromPrivate(prk)
	if err != nil {
		return err
	}

	pukder, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return err
	}
	if err = p.setPublicKey(&pem.Block{
		Type:  PublicKey.PEMString(),
		Bytes: pukder,
	}); err != nil {
		return err
	}
	return nil

}

func (p *PrivateKeyDTO) setPublicKey(block *pem.Block) error {
	if p.PublicKey != "" {
		logger.Warning("multiple public keys found in private key resource.")
		return nil
	}
	p.PublicKey = string(pem.EncodeToMemory(block))
	id := Identity(p.PublicKey)

	puk, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	p.PublicKeyAlgorithm = utils.PublicKeyAlgorithmFromKey(puk).String()
	p.Id = id.String()
	return nil
}
