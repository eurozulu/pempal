package pemwriter

import (
	"encoding/pem"
	"fmt"
	"io"
)

type PemWriter interface {
	Write(b *pem.Block) error
}

func NewPemWriter(format string, out io.Writer) (PemWriter, error) {
	switch format {
	case "", "template":
		return NewTemplateFormat(out), nil
	case "pem", "PEM":
		return NewPemFormat(out), nil
	case "der", "DER":
		return NewDerWriter(out), nil
	default:
		return nil, fmt.Errorf("%s is not a supported output format", format)
	}
}
