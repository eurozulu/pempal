package fileformats

import (
	"encoding/pem"
	"software.sslmate.com/src/go-pkcs12"
)

type pkcs12Format struct {
	Password string
}

func (p pkcs12Format) Unmarshal(by []byte) ([]*pem.Block, error) {
	return pkcs12.ToPEM(by, p.Password)
}
