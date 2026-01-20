package resourceformat

import (
	"encoding/json"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
	"io"
)

type JsonFormat struct{}

func (y JsonFormat) Format(out io.Writer, p *model.PemFile) error {
	if len(p.Blocks) == 0 {
		return nil
	}
	temps, err := templates.TemplatesOfResources(p.Resources())
	if err != nil {
		return err
	}
	var v interface{} = temps
	if len(temps) == 1 {
		v = temps[0]
	}
	return json.NewEncoder(out).Encode(v)
}
