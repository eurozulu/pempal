package resources

import (
	"bytes"
	"encoding/pem"
	"os"
)

type PemWriter interface {
	WritePem(blk ...*pem.Block) ([]byte, error)
}

type pemFileWriter struct {
	filename string
}

func (p pemFileWriter) WritePem(blk ...*pem.Block) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, block := range blk {
		if err := pem.Encode(buf, block); err != nil {
			return nil, err
		}
	}
	perm := os.FileMode(0600)
	//if !containsPrivateKey(blk) {
	perm = 0644
	//}
	if err := os.WriteFile(p.filename, buf.Bytes(), perm); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewPemFileWriter(filename string) PemWriter {
	return &pemFileWriter{filename: filename}
}
