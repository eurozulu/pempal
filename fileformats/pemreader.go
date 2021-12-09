package fileformats

import (
	"bytes"
	"encoding/pem"
	"fmt"
)

type pemReader struct {
}

func (pr pemReader) Unmarshal(by []byte) ([]*pem.Block, error) {
	if !bytes.Contains(by, []byte("-----BEGIN")) {
		return nil, fmt.Errorf("not a pem")
	}
	var blocks []*pem.Block
	for len(by) > 0 {
		block, buf := pem.Decode(by)
		if block == nil {
			break
		}
		blocks = append(blocks, block)
		by = buf
	}
	return blocks, nil
}
