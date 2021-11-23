package parsers

import (
	"encoding/pem"
)

type UnknownParser struct {
}

func (u UnknownParser) Parse(b *pem.Block) (map[string]interface{}, error) {
	return map[string]interface{}{
		PEM_TYPE: b.Type,
	}, nil
}
