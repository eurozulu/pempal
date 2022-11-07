package key

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"pempal/formats/formathelpers"
	"pempal/resources"
	"pempal/templates"
)

const encryptAlg = x509.PEMCipherAES256

type keyBuilder struct {
	Password []byte
	keyTemp  templates.KeyTemplate
	location string
}

func (fm *keyBuilder) SetLocation(l string) {
	fm.location = l
}

func (fm *keyBuilder) AddTemplate(ts ...templates.Template) error {
	for _, t := range ts {
		by, err := yaml.Marshal(t)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(by, &fm.keyTemp); err != nil {
			return err
		}
	}
	return nil
}

func (fm keyBuilder) Template() templates.Template {
	return &fm.keyTemp
}

func (fm keyBuilder) Build() (resources.Resources, error) {
	if fm.keyTemp.KeyType == "" {
		return nil, fmt.Errorf("missing key-type")
	}
	pka := formathelpers.ParsePublicKeyAlgorithm(fm.keyTemp.KeyType)
	if pka == x509.UnknownPublicKeyAlgorithm {
		return nil, fmt.Errorf("key-type unknown")
	}
	if fm.keyTemp.IsEncrypted && len(fm.Password) == 0 {
		return nil, fmt.Errorf("missing password")
	}

	var res resources.Resources
	der, err := MakeKey(pka, fm.keyTemp.Size)
	if err != nil {
		return nil, err
	}
	derPuk, err := publicBytesFromPrivate(der)
	res = append(res, resources.NewResource(fm.location,
		&pem.Block{
			Type:  resources.PublicKey.String(),
			Bytes: derPuk,
		}))
	var blk *pem.Block
	if fm.keyTemp.IsEncrypted {
		blk, err = x509.EncryptPEMBlock(rand.Reader,
			resources.PrivateKeyEncrypted.String(),
			der,
			fm.Password,
			encryptAlg)
		if err != nil {
			return nil, err
		}
	} else {
		// unencrypted key
		blk = &pem.Block{
			Type:  resources.PrivateKey.String(),
			Bytes: der,
		}
	}
	res = append(res, resources.NewResource(fm.location, blk))
	return res, nil
}

func publicBytesFromPrivate(der []byte) ([]byte, error) {
	prk, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, err
	}
	puk := PublicKeyFromPrivate(prk)
	return x509.MarshalPKIXPublicKey(puk)
}

func NewKeyFormat() *keyBuilder {
	return &keyBuilder{}
}
