package builders

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"pempal/model"
	"pempal/templates"
	"pempal/utils"
)

type PrivateKeyBuilder struct {
	dto    model.PrivateKeyDTO
	passwd []byte
}

func (kb *PrivateKeyBuilder) ApplyTemplate(tp ...templates.Template) error {
	for _, t := range tp {
		if err := t.Apply(&kb.dto); err != nil {
			return err
		}
	}
	return nil
}

func (kb PrivateKeyBuilder) Validate() []error {
	var errs []error
	m := kb.RequiredValues()
	for k := range m {
		errs = append(errs, fmt.Errorf("%s invalid", k))
	}

	if kb.dto.IsEncrypted && len(kb.passwd) == 0 {
		errs = append(errs, fmt.Errorf("No password provided to encrypt new key"))
	}
	return errs
}

func (kb PrivateKeyBuilder) RequiredValues() map[string]interface{} {
	m := map[string]interface{}{}
	if kb.dto.PublicKeyAlgorithm == "" {
		m["public-key-algorithm"] = x509.UnknownPublicKeyAlgorithm
	}
	return m
}

func (kb PrivateKeyBuilder) Build() (model.PEMResource, error) {
	if errs := kb.Validate(); len(errs) > 0 {
		return nil, fmt.Errorf("%s", collectErrorList(errs))
	}
	keyAlgo := utils.ParsePublicKeyAlgorithm(kb.dto.PublicKeyAlgorithm)

	key, err := utils.CreatePrivateKey(keyAlgo, kb.dto.KeyParam)
	if err != nil {
		return nil, err
	}

	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	var blk *pem.Block
	if !kb.dto.IsEncrypted {
		blk = &pem.Block{
			Type:  model.PrivateKey.PEMString(),
			Bytes: der,
		}
	} else {
		newDto := &model.PrivateKeyDTO{}
		if err := newDto.UnmarshalBinary(der); err != nil {
			return nil, err
		}
		if err := newDto.Encrypt(kb.passwd); err != nil {
			return nil, err
		}
		blk, err = newDto.ToPEMBlock()
		if err != nil {
			return nil, err
		}
	}
	return model.NewPemResourceFromBlock(blk), nil
}
