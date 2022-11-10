package encoders

import (
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"pempal/pemtypes"
	"pempal/templates"
)

type PublicKeyDecoder struct {
}

func (p PublicKeyDecoder) Decode(t templates.Template) (*pem.Block, error) {
	nt, ok := t.(*templates.PublicKeyTemplate)
	if !ok {
		return nil, fmt.Errorf("template is not for a Public Key")
	}
	if nt.PublicKey == "" {
		return nil, fmt.Errorf("no public key present")
	}
	data, err := hex.DecodeString(nt.PublicKey)
	if err != nil {
		return nil, err
	}
	return &pem.Block{
		Type:  pemtypes.PublicKey.String(),
		Bytes: data,
	}, nil
}
