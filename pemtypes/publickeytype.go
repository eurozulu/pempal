package pemtypes

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"pempal/misc"
	"pempal/templates"
)

type publicKeyType struct {
	puk crypto.PublicKey
}

func (pt publicKeyType) String() string {
	t := &templates.PublicKeyTemplate{}
	pt.applyToTemplate(t)
	return fmt.Sprintf("%s\t%s", t.KeyType, t.Size)
}

func (pt publicKeyType) MarshalBinary() (data []byte, err error) {
	k, ok := pt.puk.(*rsa.PublicKey)
	if ok {
		return x509.MarshalPKCS1PublicKey(k), nil
	}
	return x509.MarshalPKIXPublicKey(pt.puk)
}

func (pt *publicKeyType) UnmarshalBinary(data []byte) error {
	k, err := x509.ParsePKCS1PublicKey(data)
	if err != nil {
		k, err = x509.ParsePKCS1PublicKey(data)
		if err != nil {
			return err
		}
	}
	pt.puk = k
	return nil
}

func (pt publicKeyType) MarshalText() (text []byte, err error) {
	if pt.puk == nil {
		return nil, nil
	}
	der, err := pt.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  PublicKey.String(),
		Bytes: der,
	}), nil
}

func (pt *publicKeyType) UnmarshalText(text []byte) error {
	blocks := ReadPEMBlocks(text, PublicKey)
	if len(blocks) == 0 {
		return fmt.Errorf("no public key pems found")
	}
	return pt.UnmarshalBinary(blocks[0].Bytes)
}

func (pt publicKeyType) MarshalYAML() (interface{}, error) {
	t := templates.PublicKeyTemplate{}
	if err := pt.applyToTemplate(&t); err != nil {
		return nil, err
	}
	return yaml.Marshal(&t)
}

func (pt *publicKeyType) UnmarshalYAML(value *yaml.Node) error {
	t := templates.PublicKeyTemplate{}
	if err := value.Decode(&t); err != nil {
		return err
	}
	return pt.applyFromTemplate(&t)
}

func (pt publicKeyType) applyToTemplate(t *templates.PublicKeyTemplate) error {
	if pt.puk == nil {
		return fmt.Errorf("public key is empty")
	}
	der, err := pt.MarshalBinary()
	if err != nil {
		return err
	}
	t.KeyType = misc.PublicKeyAlgorithmFromKey(pt.puk).String()
	t.Size = misc.SizeFromKey(pt.puk)
	t.PublicKey = string(der)
	return nil
}

func (pt *publicKeyType) applyFromTemplate(t *templates.PublicKeyTemplate) error {
	if t.PublicKey != "" {
		return pt.UnmarshalBinary([]byte(t.PublicKey))
	}
	return nil
}
