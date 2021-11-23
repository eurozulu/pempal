package fileformats

import (
	"bytes"
	"encoding/pem"
	"fmt"
)

type pemFileFormat struct{}

func (p pemFileFormat) Format(by []byte) ([]*pem.Block, error) {
	if !bytes.Contains(by, []byte("---BEGIN")) {
		return nil, fmt.Errorf("not a pem format")
	}
	var blks []*pem.Block
	data := by
	for len(data) > 0 {
		blk, rest := pem.Decode(data)
		if blk == nil {
			break
		}
		blks = append(blks, blk)
		data = rest
	}
	return blks, nil
}
