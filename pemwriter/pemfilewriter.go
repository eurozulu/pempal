package pemwriter

import (
	"encoding/pem"
	"io"
)

type pemWriter struct {
	out io.Writer
}

func (t pemWriter) Write(b *pem.Block) error {
	return pem.Encode(t.out, b)
}

func NewPemFormat(out io.Writer) *pemWriter {
	return &pemWriter{
		out: out,
	}
}
