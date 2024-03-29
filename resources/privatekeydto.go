package resources

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/utils"
	"strings"
)

var errNoPrivateKey = fmt.Errorf("no pem encoded private key found")

type PrivateKeyDTO struct {
	// PublicKeyAlgorithm is the Key Algorithm of the existing key or key to be created
	PublicKeyAlgorithm string `yaml:"key-algorithm" json:"key-algorithm"`
	IsEncrypted        bool   `yaml:"is-encrypted" json:"is-encrypted"`
	KeyLength          string `yaml:"key-length,omitempty" json:"key-length"`
	KeyCurve           string `yaml:"key-curve,omitempty" json:"key-curve"`
	PrivateKey         string `yaml:"private-key,omitempty" json:"private-key,omitempty"`
	// PublicKey when it exists, contains the PEM encoded public key.
	PublicKey string `yaml:"public-key,omitempty" json:"-"`
}

func (pk PrivateKeyDTO) String() string {
	return strings.Join([]string{pk.PrivateKey, pk.PublicKey}, "")
}

func (pk *PrivateKeyDTO) UnmarshalPEM(data []byte) error {
	blk, _ := pem.Decode(data)
	if blk == nil {
		return fmt.Errorf("failed to parse private key from pem")
	}
	if x509.IsEncryptedPEMBlock(blk) {
		return pk.unmarshalEncryptedPEM(data)
	}
	return pk.UnmarshalBinary(blk.Bytes)
}

func (pk *PrivateKeyDTO) unmarshalEncryptedPEM(data []byte) error {
	pk.PrivateKey = string(data)
	pk.PublicKey = ""
	pk.IsEncrypted = true
	return nil
}

func (pk *PrivateKeyDTO) UnmarshalBinary(data []byte) error {
	prk, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return err
	}
	pk.PrivateKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  PrivateKey.PEMString(),
		Bytes: data,
	}))
	puk, err := utils.PublicKeyFromPrivate(prk)
	if err != nil {
		return err
	}
	pk.PublicKeyAlgorithm = utils.PublicKeyAlgorithmFromKey(puk).String()

	if pukPem, err := utils.PublicKeyToPEM(puk); err == nil {
		pk.PublicKey = string(pem.EncodeToMemory(pukPem))
	} else {
		pk.PublicKey = ""
		logger.Warning("failed to encode public key to pem %v", err)
	}
	pk.IsEncrypted = false
	return nil
}

func (pk PrivateKeyDTO) MarshalBinary() (data []byte, err error) {
	blk, _ := pem.Decode([]byte(pk.PrivateKey))
	if blk == nil {
		return nil, errNoPrivateKey
	}
	return blk.Bytes, nil

}
