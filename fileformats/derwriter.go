package fileformats

import (
	"encoding/pem"
	"io"
)

const defaultDelimiter = "\n"

type derWriter struct {
	Delimiter []byte
}

func (d derWriter) Marshal(blks []*pem.Block, out io.Writer) error {
	dl := d.Delimiter
	if len(dl) == 0 {
		dl = []byte(defaultDelimiter)
	}

	for i, blk := range blks {
		if i > 0 {
			if _, err := out.Write(dl); err != nil {
				return err
			}
			if _, err := out.Write(blk.Bytes); err != nil {
				return err
			}
		}
	}
	return nil
}
