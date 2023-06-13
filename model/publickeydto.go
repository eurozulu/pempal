package model

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/utils"
)

type PublicKeyDTO struct {
	Id                 string `yaml:"identity,omitempty"`
	PublicKeyAlgorithm string `yaml:"public-key-algorithm"`
	PublicKey          string `yaml:"public-key,omitempty"`

	ResourceType string `yaml:"resource-type"`
}

func (p *PublicKeyDTO) UnmarshalPEM(data []byte) error {
	for len(data) > 0 {
		blk, rest := pem.Decode(data)
		if blk == nil {
			break
		}
		if ParsePEMType(blk.Type) != PublicKey {
			data = rest
			continue
		}
		return p.UnmarshalBinary(blk.Bytes)
	}
	return fmt.Errorf("no pem encoded public key found")
}

func (p PublicKeyDTO) String() string {
	return p.PublicKey
}

func (p *PublicKeyDTO) MarshalBinary() (data []byte, err error) {
	blk, _ := pem.Decode([]byte(p.PublicKey))
	if blk == nil {
		return nil, fmt.Errorf("failed to marshal public key.  Invalid PublicKey property")
	}
	return blk.Bytes, nil
}

func (pkr *PublicKeyDTO) UnmarshalBinary(data []byte) error {
	pkr.PublicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  PublicKey.PEMString(),
		Bytes: data,
	}))
	pkr.Id = Identity(pkr.PublicKey).String()
	puk, err := pkr.ToPublicKey()
	if err != nil {
		return err
	}
	pkr.PublicKeyAlgorithm = utils.PublicKeyAlgorithmFromKey(puk).String()
	pkr.ResourceType = PublicKey.String()
	return nil
}

func (p PublicKeyDTO) ToPublicKey() (crypto.PublicKey, error) {
	der, err := p.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("public-key is invalid hex or empty %v", err)
	}
	return x509.ParsePKIXPublicKey(der)
}

func NewPublicKeyDTO(puk crypto.PublicKey) (*PublicKeyDTO, error) {
	der, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return nil, err
	}
	dto := &PublicKeyDTO{}
	if err = dto.UnmarshalBinary(der); err != nil {
		return nil, err
	}
	return dto, nil
}
