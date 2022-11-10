package encoders

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"pempal/templates"
)

type PrivateKeyEncoder struct {
}

func (pke PrivateKeyEncoder) Encode(p *pem.Block) (templates.Template, error) {
	if x509.IsEncryptedPEMBlock(p) {
		t := &templates.KeyTemplate{}
		t.IsEncrypted = true
		return t, nil
	}

	prk, err := ParsePrivateKey(p.Bytes)
	if err != nil {
		prk, err = x509.ParsePKCS1PrivateKey(p.Bytes)
		if err != nil {
			return nil, err
		}
	}
	t := &templates.KeyTemplate{}
	pke.ApplyPem(prk, t)
	return t, nil
}

func (pke PrivateKeyEncoder) ApplyPem(prk crypto.PrivateKey, t *templates.KeyTemplate) {
	puk := PublicKeyFromPrivate(prk)
	t.PublicKey = &templates.PublicKeyTemplate{}
	PublicKeyEncoder{}.ApplyPem(puk, t.PublicKey)
}
