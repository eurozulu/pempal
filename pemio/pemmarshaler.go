package pemio

import (
	"fmt"
	"github.com/pempal/templates"
	"gopkg.in/yaml.v3"
	"io"

	"encoding/pem"
)

type PEMMarshaler interface {
	Marshal(out io.Writer, pb []*pem.Block) error
}

func NewPEMMarshaler(format string) PEMMarshaler {
	switch format {
	case "pem":
		return &pemMarshaler{}
	case "der":
		return &derMarshaler{}
	case "p12":
		return &p12Marshaler{}
	case "yaml", "text":
		return &yamlMarshaler{}
	default:
		return nil
	}
}

type pemMarshaler struct{}

func (pm pemMarshaler) Marshal(out io.Writer, blks []*pem.Block) error {
	for _, b := range blks {
		if err := pem.Encode(out, b); err != nil {
			return err
		}
	}
	return nil
}

type derMarshaler struct{}

func (pm derMarshaler) Marshal(out io.Writer, blks []*pem.Block) error {
	for _, b := range blks {
		if _, err := out.Write(b.Bytes); err != nil {
			return err
		}
	}
	return nil
}

type p12Marshaler struct{}

func (pm p12Marshaler) Marshal(out io.Writer, blks []*pem.Block) error {
	panic("Not yet implemented")
}

type yamlMarshaler struct{}

func (pm yamlMarshaler) Marshal(out io.Writer, blks []*pem.Block) error {
	enc := yaml.NewEncoder(out)
	for _, bl := range blks {
		t, err := templates.NewTemplate(bl)
		if err != nil {
			return err
		}
		if err := enc.Encode(t); err != nil {
			return fmt.Errorf("faileds to encode PEM into yaml  %v", err)
		}
	}
	return nil
}
