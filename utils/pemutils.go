package utils

import "encoding/pem"

func ReadPEMBlocks(data []byte) []*pem.Block {
	var blocks []*pem.Block
	for len(data) > 0 {
		b, rest := pem.Decode(data)
		if b == nil {
			break
		}
		blocks = append(blocks, b)
		data = rest
	}
	return blocks
}
