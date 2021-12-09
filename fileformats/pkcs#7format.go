package fileformats

import (
	"encoding/pem"
)

type pkcs7Format struct {
}

func (p pkcs7Format) Unmarshal(by []byte) ([]*pem.Block, error) {
	panic("implement me")
}
