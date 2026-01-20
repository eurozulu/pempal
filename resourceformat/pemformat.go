package resourceformat

import (
	"encoding/pem"
	"github.com/eurozulu/pempal/model"
	"io"
)

type PemFormat struct{}

func (pf PemFormat) Format(out io.Writer, p *model.PemFile) error {
	for i, block := range p.Blocks {
		if i > 0 {
			out.Write([]byte("\n"))
		}
		if err := pem.Encode(out, block); err != nil {
			return err
		}
	}
	return nil
}
