package templates

import (
	"github.com/eurozulu/pempal/model"
	"gopkg.in/yaml.v2"
)

type PrivateKeyTemplate struct {
	KeyAlgoritum  model.PublicKeyAlgorithm `yaml:"public-key-algorithm"`
	RSAKeyLength  int                      `yaml:"key-length,omitempty"`
	ECDSAKeyCurve model.EllipticCurve      `yaml:"key-curve,omitempty"`
}

func (t PrivateKeyTemplate) String() string {
	data, err := yaml.Marshal(&t)
	if err != nil {
		return ""
	}
	return string(data)
}

func NewPrivateKeyTemplate(r *model.PrivateKey) *PrivateKeyTemplate {
	return &PrivateKeyTemplate{
		KeyAlgoritum:  r.PublicKeyAlgorithm(),
		RSAKeyLength:  r.RSAKeyLength(),
		ECDSAKeyCurve: r.ECDSACurve(),
	}
}
