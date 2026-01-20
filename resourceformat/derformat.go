package resourceformat

import (
	"encoding"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"io"
)

type DERFormat struct{}

func (D DERFormat) Format(out io.Writer, p *model.PemFile) error {
	if len(p.Blocks) == 0 {
		return nil
	}
	if len(p.Blocks) > 1 {
		return fmt.Errorf("multiple PEM resources in %s.  derformat can only format a single block", p.Path)
	}
	res := p.Resources()[0]
	r, ok := res.(encoding.BinaryMarshaler)
	if !ok {
		return fmt.Errorf("resource type %s does not support DER encoding", res.ResourceType().String())
	}
	data, err := r.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = out.Write(data)
	return err
}
