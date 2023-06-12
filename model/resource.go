package model

import (
	"encoding/pem"
	"fmt"
)

// Resource represents a single PEM resource.
// The resource may be any one of the ResourceType resources
type Resource interface {
	// ResourceType is the type of the resource based on Its PEM type.
	ResourceType() ResourceType

	// String returns the resource in its PEM encoded format
	fmt.Stringer
}

func (r resource) String() string {
	return string(pem.EncodeToMemory(r.block))
}

type resource struct {
	block *pem.Block
}

func (r resource) ResourceType() ResourceType {
	return ParsePEMType(r.block.Type)
}

func NewResource(block *pem.Block) Resource {
	return &resource{block: block}
}
