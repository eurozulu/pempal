package model

import (
	"encoding"
	"encoding/pem"
	"fmt"
)

type PEMResource interface {
	ResourceType() ResourceType
	// BinaryMarshaler and BinaryUnmarshaler encode/decode into der encoded binary
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	MarshalPEM() (data []byte, err error)
	UnmarshalPEM(data []byte) error
}

type pemResource struct {
	block *pem.Block
}

func (r pemResource) ResourceType() ResourceType {
	return ParsePEMType(r.block.Type)
}

func (r pemResource) MarshalBinary() (data []byte, err error) {
	if r.block == nil {
		return nil, fmt.Errorf("pemResource is empty")
	}
	return r.block.Bytes, nil
}

func (r *pemResource) UnmarshalBinary(data []byte) error {
	if r.block == nil {
		r.block = &pem.Block{
			Type: Unknown.PEMString(),
		}
	}
	r.block.Bytes = data
	return nil
}

func (r pemResource) MarshalPEM() (data []byte, err error) {
	if r.block == nil {
		return nil, fmt.Errorf("pemResource is empty")
	}
	return pem.EncodeToMemory(r.block), nil
}

func (r *pemResource) UnmarshalPEM(data []byte) error {
	r.block, _ = pem.Decode(data)
	if r.block == nil {
		return fmt.Errorf("failed to decode as pem")
	}
	return nil
}

func NewPemResource() PEMResource {
	return &pemResource{}
}

func NewPemResourceFromBlock(pemblock *pem.Block) PEMResource {
	return &pemResource{block: pemblock}
}
