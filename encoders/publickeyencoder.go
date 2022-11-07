package encoders

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"pempal/templates"
)

type publicKeyEncoder struct{}

func (pe publicKeyEncoder) Encode(p *pem.Block) (templates.Template, error) {
	puk, err := x509.ParsePKIXPublicKey(p.Bytes)
	if err != nil {
		return nil, err
	}
	t := &templates.PublicKeyTemplate{}
	pe.ApplyPem(puk, t)
	return t, nil
}

func (pe publicKeyEncoder) ApplyPem(puk crypto.PublicKey, t *templates.PublicKeyTemplate) {
	t.PublicKey = MarshalPublicKey(puk)
	t.Size = MarshalSizeFromKey(puk)
	t.KeyType = PublicKeyAlgorithmFromKey(puk).String()
}
