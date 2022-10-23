package resources

import (
	"encoding/pem"
	"fmt"
)

type pemResource struct {
	block    *pem.Block
	location string
}

func (c pemResource) ResourceType() string {
	return c.block.Type
}

func (c pemResource) Location() string {
	return c.location
}

func (c pemResource) MarshalBinary() (data []byte, err error) {
	return c.MarshalPEM()
}

func (c pemResource) UnmarshalBinary(data []byte) error {
	return c.UnmarshalPem(data)
}

func (c pemResource) MarshalPEM() (data []byte, err error) {
	if c.block == nil {
		return nil, nil
	}
	return pem.EncodeToMemory(c.block), nil
}

func (c *pemResource) UnmarshalPem(data []byte) error {
	// read first pem block
	blk, _ := pem.Decode(data)
	if blk == nil {
		return fmt.Errorf("no pem found")
	}
	c.block = blk
	return nil
}
