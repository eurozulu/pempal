package model

import (
	"bytes"
	"encoding/pem"
	"fmt"
)

func parseAsBinary(data []byte) ([]PEMResource, error) {
	for _, rt := range resourceTypes {
		d := DTOForResourceType(rt)
		if err := d.UnmarshalBinary(data); err != nil {
			continue
		}
		return []PEMResource{&pemResource{block: &pem.Block{
			Type:  rt.PEMString(),
			Bytes: data,
		}}}, nil
	}
	return nil, fmt.Errorf("Failed to parse as der resources")
}

func parseAsPem(data []byte) ([]PEMResource, error) {
	var res []PEMResource
	blocks := readPEMBlocks(data)
	for _, blk := range blocks {
		res = append(res, &pemResource{block: blk})
	}
	return res, nil
}

func readPEMBlocks(data []byte, types ...ResourceType) []*pem.Block {
	var blocks []*pem.Block
	for len(data) > 0 {
		b, rest := pem.Decode(data)
		if b == nil {
			break
		}
		if len(types) > 0 && !ContainsType(ParsePEMType(b.Type), types) {
			continue
		}
		blocks = append(blocks, b)
		data = rest
	}
	return blocks
}

func isPEMBytes(data []byte) bool {
	return bytes.Contains(data, []byte("-----BEGIN ")) && bytes.Contains(data, []byte("-----END "))
}

func ParseResources(data []byte) ([]PEMResource, error) {
	if isPEMBytes(data) {
		return parseAsPem(data)
	}
	return parseAsBinary(data)
}
