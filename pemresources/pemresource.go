package pemresources

import (
	"encoding/pem"
)

// PemResource represents a generic pem resource.
// It contains all the properties of the pem, excluding the 'bytes', with the addition of anoptional Location
type PemResource struct {
	//PemType is the pem type, CERTIFICATE, PRIVATE KEY etc
	PemType string `yaml:"pem_type"`
	// PemHeaders are optional headers found in a pem block
	PemHeaders map[string]string `yaml:"pemHeaders,omitempty"`
	// Location is an optional filepath of the location of the pem.
	Location string `yaml:"location,omitempty"`
}

func (pt *PemResource) MarshalPem() (*pem.Block, error) {
	return &pem.Block{
		Type:    pt.PemType,
		Headers: pt.PemHeaders,
		Bytes:   nil,
	}, nil
}

func (pt *PemResource) UnmarshalPem(block *pem.Block) error {
	pt.PemType = block.Type
	pt.PemHeaders = block.Headers
	return nil
}
func (pt PemResource) Type() string {
	return pt.PemType
}
