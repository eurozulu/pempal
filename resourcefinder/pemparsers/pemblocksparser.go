package pemparsers

import (
	"encoding/pem"
	"path/filepath"
	"strings"
)

var pemFileExtensions = []string{"", "pem", "cer", "crt", "key", "csr", "cert", "crl"}

type pemBlocksParser struct {
}

func (p pemBlocksParser) Match(path string) bool {
	return stringIndex(strings.ToLower(filepath.Ext(path)), pemFileExtensions) >= 0

}

func (p pemBlocksParser) Parse(data []byte) (PemBlocks, []byte, error) {
	var res PemBlocks
	var blk *pem.Block
	for len(data) > 0 {
		blk, data = pem.Decode(data)
		if blk == nil {
			break
		}
		res = append(res, blk)
	}
	return res, data, nil
}
