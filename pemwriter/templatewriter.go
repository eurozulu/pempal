package pemwriter

import (
	"encoding/pem"
	"gopkg.in/yaml.v3"
	"io"
	"pempal/templates"
)

type templateWriter struct {
	out     io.Writer
	encoder *yaml.Encoder
}

func (tw templateWriter) Write(b *pem.Block) error {
	t, err := templates.ParseBlock(b)
	if err != nil {
		return err
	}
	return tw.encoder.Encode(t)
}

func NewTemplateFormat(out io.Writer) *templateWriter {
	return &templateWriter{
		out:     out,
		encoder: yaml.NewEncoder(out),
	}
}
