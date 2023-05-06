package model

import (
	"crypto"
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"pempal/utils"
)

type PublicKeyDTO struct {
	PublicKeyAlgorithm string `yaml:"public-key-algorithm"`
	PublicKey          string `yaml:"public-key,omitempty"`

	Identity     string `yaml:"identity"`
	ResourceType string `yaml:"resource-type"`
}

func (p PublicKeyDTO) String() string {
	if p.PublicKey == "" {
		return ""
	}
	der, err := hex.DecodeString(p.PublicKey)
	if err != nil {
		return ""
	}
	return MD5PublicKey(der)
}

func (p PublicKeyDTO) ToPublicKey() (crypto.PublicKey, error) {
	der, err := p.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("public-key is invalid hex or empty %v", err)
	}
	return x509.ParsePKIXPublicKey(der)
}

func (p PublicKeyDTO) ToPem() ([]byte, error) {
	der, err := p.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("public-key is invalid hex or empty %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  PublicKey.PEMString(),
		Bytes: der,
	}), nil
}

func (p *PublicKeyDTO) MarshalBinary() (data []byte, err error) {
	return hex.DecodeString(p.PublicKey)
}

func (pkr *PublicKeyDTO) UnmarshalBinary(data []byte) error {
	pkr.PublicKey = hex.EncodeToString(data)
	puk, err := pkr.ToPublicKey()
	if err != nil {
		return err
	}
	pkr.PublicKeyAlgorithm = utils.PublicKeyAlgorithmFromKey(puk).String()
	pkr.Identity = pkr.String()
	pkr.ResourceType = PublicKey.String()
	return nil
}

func NewPublicKeyDTO(puk crypto.PublicKey) (PublicKeyDTO, error) {
	keydto := PublicKeyDTO{}
	der, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return keydto, err
	}
	return PublicKeyDTO{
		PublicKeyAlgorithm: utils.PublicKeyAlgorithmFromKey(puk).String(),
		PublicKey:          hex.EncodeToString(der),
	}, err
}

func MD5PublicKey(der []byte) string {
	hash := md5.Sum(der)
	return hex.EncodeToString(hash[:])
}
