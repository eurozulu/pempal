package resourceformat

import (
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v2"
	"io"
)

type YamlFormat struct{}

func (y YamlFormat) Format(out io.Writer, p *model.PemFile) error {
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
	return yaml.NewEncoder(out).Encode(v)
}
