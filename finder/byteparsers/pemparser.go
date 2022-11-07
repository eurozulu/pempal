package byteparsers

import (
	"encoding/pem"
	"path/filepath"
	"strings"
)

var pemFileExtensions = map[string]bool{"": true, "pem": true, "cer": true, "crt": true, "key": true, "csr": true, "cert": true, "crl": true}

type pemParser struct {
}

func (p pemParser) MatchPath(path string) bool {
	return pemFileExtensions[strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))]

}

func (p pemParser) Parse(data []byte) ([]*pem.Block, []byte, error) {
	var res []*pem.Block
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
