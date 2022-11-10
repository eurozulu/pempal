package encoders

import (
	"crypto"
	"encoding/pem"
	"pempal/templates"
)

type PublicKeyEncoder struct{}

func (pe PublicKeyEncoder) Encode(p *pem.Block) (templates.Template, error) {
	puk, err := ParsePublicKey(p.Bytes)
	if err != nil {
		return nil, err
	}
	t := &templates.PublicKeyTemplate{}
	pe.ApplyPem(puk, t)
	return t, nil
}

func (pe PublicKeyEncoder) ApplyPem(puk crypto.PublicKey, t *templates.PublicKeyTemplate) {
	t.PublicKey = MarshalPublicKey(puk)
	t.Size = MarshalSizeFromKey(puk)
	t.KeyType = PublicKeyAlgorithmFromKey(puk).String()
}
