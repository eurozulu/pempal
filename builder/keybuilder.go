package builder

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"github.com/go-yaml/yaml"
)

type keyBuilder struct {
	dto model.PrivateKeyDTO
}

func (kb *keyBuilder) ApplyTemplate(tp ...templates.Template) error {
	for _, t := range tp {
		if err := yaml.Unmarshal(t.Bytes(), &kb.dto); err != nil {
			return err
		}
	}
	return nil
}

func (kb keyBuilder) Validate() error {
	if _, err := kb.Build(); err != nil {
		return err
	}
	return nil
}

func (kb keyBuilder) Build() (model.Resource, error) {
	errs := bytes.NewBuffer(nil)
	pka := utils.ParsePublicKeyAlgorithm(kb.dto.PublicKeyAlgorithm)
	if pka == x509.UnknownPublicKeyAlgorithm {
		errs.WriteString("PublicKey Algorithm Unknown")
	}
	length := kb.dto.KeySize
	if length == "" {
		errs.WriteString("Keysize Unknown")
	}
	prk, err := utils.CreatePrivateKey(pka, length)
	if err != nil {
		return nil, err
	}
	der, err := x509.MarshalPKCS8PrivateKey(prk)
	if err != nil {
		return nil, err
	}
	return model.NewResource(&pem.Block{
		Type:  model.PrivateKey.PEMString(),
		Bytes: der,
	}), nil
}
