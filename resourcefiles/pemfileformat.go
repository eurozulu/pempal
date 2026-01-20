package resourcefiles

import (
	"encoding/pem"
)

type PemFileFormat struct {
}

func (p PemFileFormat) Format(data []byte) ([]*pem.Block, error) {
	var found []*pem.Block
	for len(data) > 0 {
		blk, rest := pem.Decode(data)
		if blk == nil {
			break
		}
		found = append(found, blk)
		data = rest
	}
	return found, nil
}
