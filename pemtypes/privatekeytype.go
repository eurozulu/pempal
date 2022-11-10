package pemtypes

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
	"pempal/misc"
	"pempal/templates"
)

type privateKeyType struct {
	prk       crypto.PrivateKey
	puk       crypto.PublicKey
	encrypted *pem.Block
}

func (pt privateKeyType) String() string {
	kt := PrivateKey
	if pt.encrypted != nil {
		kt = PrivateKeyEncrypted
	}
	var pukS string
	if pt.puk != nil {
		t := pt.publicKeyTemplate()
		pukS = fmt.Sprintf("%s\t%s", t.KeyType, t.Size)
	}
	return fmt.Sprintf("%s\t%s", kt.String(), pukS)
}

func (pt privateKeyType) MarshalBinary() (data []byte, err error) {
	k, ok := pt.prk.(*rsa.PrivateKey)
	if ok {
		return x509.MarshalPKCS1PrivateKey(k), nil
	}
	return x509.MarshalPKCS8PrivateKey(pt.prk)
}

func (pt privateKeyType) UnmarshalBinary(data []byte) error {
	k, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		k, err = x509.ParsePKCS1PrivateKey(data)
		if err != nil {
			return err
		}
	}
	pt.prk = k
	pt.puk = misc.PublicKeyFromPrivate(k)
	return nil
}

// MarshalText marshals this private key into a PEM encoded blocks.
// Outputs two blocks (when present), both the private and public keys.
// If key is encrypted, public key, if present and the encrypted key
func (pt privateKeyType) MarshalText() (text []byte, err error) {
	if pt.prk == nil && pt.encrypted == nil {
		return nil, nil
	}
	buf := bytes.NewBuffer(nil)
	var p *pem.Block
	if pt.encrypted == nil {
		der, err := pt.MarshalBinary()
		if err != nil {
			return nil, err
		}
		p = &pem.Block{
			Type:  PrivateKey.String(),
			Bytes: der,
		}
	} else {
		p = pt.encrypted
	}
	pem.Encode(buf, p)

	if pt.puk != nil {
		put := &publicKeyType{puk: pt.puk}
		der, err := put.MarshalBinary()
		if err != nil {
			return nil, err
		}
		pem.Encode(buf, &pem.Block{
			Type:  PublicKey.String(),
			Bytes: der,
		})
	}
	return buf.Bytes(), nil
}

func (pt privateKeyType) UnmarshalText(text []byte) error {
	privateKeyBlocks := ReadPEMBlocks(text, PrivateKey, PrivateKeyEncrypted)
	if len(privateKeyBlocks) == 0 {
		return fmt.Errorf("no private key pems found")
	}
	ty := ParsePEMType(privateKeyBlocks[0].Type)
	if ty == PrivateKeyEncrypted || x509.IsEncryptedPEMBlock(privateKeyBlocks[0]) {
		pt.encrypted = privateKeyBlocks[0]
		// If encrypted, look for a public key in the same data
		pkt := &publicKeyType{}
		if err := pkt.UnmarshalBinary(text); err == nil {
			pt.puk = pkt.puk
		}
		return nil
	}
	return pt.UnmarshalBinary(privateKeyBlocks[0].Bytes)
}

func (pt privateKeyType) MarshalYAML() (interface{}, error) {
	t := templates.KeyTemplate{}
	if pt.puk != nil {
		t.PublicKey = pt.publicKeyTemplate()
	}
	t.IsEncrypted = pt.encrypted != nil
	return yaml.Marshal(&t)
}

func (pt privateKeyType) UnmarshalYAML(value *yaml.Node) error {
	t := templates.KeyTemplate{}
	if err := value.Decode(&t); err != nil {
		return err
	}
	pt.prk = nil
	pt.puk = nil
	pt.encrypted = nil
	if t.IsEncrypted {
		pt.encrypted = &pem.Block{}
	}

	if t.PublicKey != nil {
		data, err := yaml.Marshal(t.PublicKey)
		if err != nil {
			return err
		}
		put := &publicKeyType{}
		if err = put.UnmarshalText(data); err != nil {
			return err
		}
		pt.puk = put.puk
	}
	return nil
}

func (pt privateKeyType) publicKeyTemplate() *templates.PublicKeyTemplate {
	put := &publicKeyType{puk: pt.puk}
	t := &templates.PublicKeyTemplate{}
	put.applyToTemplate(t)
	return t
}
