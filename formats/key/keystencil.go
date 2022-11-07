package key

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"pempal/formats/formathelpers"
	"pempal/resources"
	"pempal/templates"
)

type keyStencil struct {
}

func (st keyStencil) MakeTemplate(r resources.Resource) (templates.Template, error) {
	puk, err := st.getPublicKey(r)
	if err != nil {
		return nil, err
	}
	t := &templates.KeyTemplate{}
	if err = st.copyToTemplate(t, puk); err != nil {
		return nil, err
	}
	return t, nil
}

func (st keyStencil) getPublicKey(r resources.Resource) (crypto.PublicKey, error) {
	blk := r.Pem()
	if blk == nil {
		return nil, fmt.Errorf("resource is empty")
	}

	switch r.Type() {
	case resources.PrivateKey:
		prk, err := x509.ParsePKCS8PrivateKey(blk.Bytes)
		if err != nil {
			return nil, err
		}
		puk := PublicKeyFromPrivate(prk)
		if puk == nil {
			return nil, fmt.Errorf("public key could not be established from the private key!!")
		}
		return puk, nil

	case resources.PublicKey:
		return x509.ParsePKIXPublicKey(blk.Bytes)

	default:
		return nil, fmt.Errorf("failed to read resource as key")
	}
}

func (st keyStencil) copyToTemplate(t *templates.KeyTemplate, puk crypto.PublicKey) error {
	t.KeyType = PublicKeyAlgorithmFromKey(puk).String()
	t.Size = formathelpers.MarshalSizeFromKey(puk)
	pk, err := formathelpers.MarshalPublicKey(puk)
	if err == nil {
		t.PublicKey = pk
	}
	return nil
}
