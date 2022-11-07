package encoders

import (
	"encoding/pem"
	"fmt"
	"pempal/pemtypes"
	"pempal/templates"
)

var encoderTypes = map[pemtypes.PEMType]Encoder{
	pemtypes.Name:        &nameEncoder{},
	pemtypes.Request:     &requestEncoder{},
	pemtypes.Certificate: &certificateEncoder{},
	pemtypes.PrivateKey:  &privateKeyEncoder{},
	pemtypes.PublicKey:   &publicKeyEncoder{},
}

type Encoder interface {
	Encode(p *pem.Block) (templates.Template, error)
}

func NewEncoder(pt pemtypes.PEMType) (Encoder, error) {
	e, ok := encoderTypes[pt]
	if !ok {
		return nil, fmt.Errorf("%s is an unknown encoder type", pt.String())
	}
	return e, nil
}
