package fileformats

import (
	"encoding/pem"
	"fmt"
	"io"
)

var FormatWriters = map[string]FormatWriter{
	"pem": &pemWriter{},
	"der": &derWriter{},
}

// FormatWriter writes pem blocks out in soecific format
type FormatWriter interface {
	Marshal(blk []*pem.Block, out io.Writer) error
}

func NewFormatWriter(format string) (FormatWriter, error) {
	fw, ok := FormatWriters[format]
	if !ok {
		return nil, fmt.Errorf("%s is not a known output format", format)
	}
	return fw, nil
}
