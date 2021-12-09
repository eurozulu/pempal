package fileformats

import (
	"encoding/pem"
	"io"
)

type pemWriter struct {
}

func (pw pemWriter) Marshal(blks []*pem.Block, out io.Writer) error {
	for _, blk := range blks {
		if err := pem.Encode(out, blk); err != nil {
			return err
		}
	}
	return nil
}
