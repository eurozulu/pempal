package pemwriter

import (
	"encoding/pem"
	"io"
)

type derWriter struct {
	out io.Writer
}

func (t derWriter) Write(b *pem.Block) error {
	buf := append(b.Bytes, '\n')
	_, err := t.out.Write(buf)
	return err
}

func NewDerWriter(out io.Writer) *derWriter {
	return &derWriter{
		out: out,
	}
}
