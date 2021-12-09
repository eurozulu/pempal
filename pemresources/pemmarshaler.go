package pemresources

import "encoding/pem"

type Marshaler interface {
	MarshalPem() (*pem.Block, error)
}

type Unmarshaler interface {
	UnmarshalPem(block *pem.Block) error
}
