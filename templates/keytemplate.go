package templates

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/utils"
	"strings"
)

type KeyTemplate struct {
	PublicKeyAlgorithm model.PublicKeyAlgorithm `yaml:"public-key-algorithm"`
	KeySize            int                      `yaml:"key-size,omitempty"`
	ID                 model.KeyId              `yaml:"id,omitempty"`
	PublicKey          model.Base64Binary       `yaml:"public-key,omitempty"`
}

func (p KeyTemplate) Name() string {
	return strings.ToLower(model.PrivateKey.String())
}

func NewPublicKeyTemplate(puk crypto.PublicKey) (*KeyTemplate, error) {
	id, err := model.NewKeyIdFromKey(puk)
	if err != nil {
		return nil, err
	}
	by, err := x509.MarshalPKIXPublicKey(puk)
	if err != nil {
		return nil, err
	}

	kt := &KeyTemplate{
		PublicKeyAlgorithm: model.PublicKeyAlgorithm(utils.PublicKeyType(puk)),
		ID:                 id,
		PublicKey:          model.Base64Binary(by),
	}
	if k, ok := puk.(*rsa.PublicKey); ok {
		kt.KeySize = k.Size()
	}
	return kt, nil
}

func NewPrivateKeyTemplate(prk crypto.PrivateKey) (*KeyTemplate, error) {
	puk, err := utils.PublicKeyFromPrivate(prk)
	if err != nil {
		return nil, err
	}
	return NewPublicKeyTemplate(puk)
}
