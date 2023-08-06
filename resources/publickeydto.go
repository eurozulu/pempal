package resources

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/utils"
)

var errNoPublicKey = fmt.Errorf("no pem encoded public key found")

type PublicKeyDTO struct {
	// PublicKeyAlgorithm is the Key Algorithm of the existing key or key to be created
	PublicKeyAlgorithm string `yaml:"public-key-algorithm" json:"public-key-algorithm"`

	// PublicKey when it exists, contains the PEM encoded public key.
	PublicKey string `yaml:"public-key,omitempty" json:"-"`
}

func (p PublicKeyDTO) String() string {
	return p.PublicKey
}

func (p *PublicKeyDTO) MarshalBinary() (data []byte, err error) {
	blk, _ := pem.Decode([]byte(p.PublicKey))
	if blk == nil {
		return nil, errNoPublicKey
	}
	return blk.Bytes, nil
}

func (pkr *PublicKeyDTO) UnmarshalBinary(data []byte) error {
	puk, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return err
	}
	pkr.PublicKeyAlgorithm = utils.PublicKeyAlgorithmFromKey(puk).String()
	pkr.PublicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  PublicKey.PEMString(),
		Bytes: data,
	}))
	return nil
}

func (pkr *PublicKeyDTO) UnmarshalPEM(data []byte) error {
	blk, _ := pem.Decode(data)
	if blk == nil {
		return fmt.Errorf("failed to parse public key from pem")
	}
	return pkr.UnmarshalBinary(blk.Bytes)
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
