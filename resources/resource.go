package resources

import (
	"encoding/pem"
)

// Resource represents a single PEM resource.
type Resource interface {
	// ResourceType is the type of the resource based on Its PEM type.
	ResourceType() ResourceType

	// String returns the resource as a string, in its PEM endcoded format.
	String() string

	// Bytes returns the raw der bytes of the resource
	Bytes() []byte
}

type resource struct {
	block *pem.Block
}

func (r resource) String() string {
	return string(pem.EncodeToMemory(r.block))
}

func (r resource) Bytes() []byte {
	if r.block == nil {
		return nil
	}
	return r.block.Bytes
}

func (r resource) ResourceType() ResourceType {
	if r.block == nil {
		return UnknownResourceType
	}
	return ParsePEMType(r.block.Type)
}

func NewResource(block *pem.Block) Resource {
	return &resource{block: block}
}
