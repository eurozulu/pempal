package encoding

import (
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/templates"
	"gopkg.in/yaml.v3"
	"io"
)

// TemplateEncoder encodes templates into a specific encoding.
type TemplateEncoder interface {
	Encode(tps []templates.Template) error
}

func NewEncoder(t string, out io.Writer) (TemplateEncoder, error) {
	switch t {
	case "pem":
		return &pemEncoder{out: out}, nil
	case "der":
		return &binEncoder{out: out}, nil
	case "yaml", "text":
		return &yamlEncoder{out: out}, nil
	default:
		return nil, fmt.Errorf("%s is an unknown encoding", t)
	}
}

type yamlEncoder struct {
	out io.Writer
}

func (f yamlEncoder) Encode(tps []templates.Template) error {
	ye := yaml.NewEncoder(f.out)
	for _, t := range tps {
		if err := ye.Encode(t); err != nil {
			return err
		}
	}
	return nil
}

type binEncoder struct {
	out io.Writer
}

func (ec binEncoder) Encode(tps []templates.Template) error {
	for _, t := range tps {
		by, err := t.MarshalBinary()
		if err != nil {
			return err
		}
		if _, err := ec.out.Write(by); err != nil {
			return err
		}
	}
	return nil
}

// pemEncoder encodes templates into pem encoded resources.
type pemEncoder struct {
	out io.Writer
}

func (pe pemEncoder) Encode(tps []templates.Template) error {
	for _, t := range tps {
		var bl *pem.Block
		if tpPem, ok := t.(PEMMarshaler); ok {
			b, err := tpPem.MarshalPEM()
			if err != nil {
				return err
			}
			bl = b
		} else {
			tp := templates.TemplateType(t)
			by, err := t.MarshalBinary()
			if err != nil {
				return err
			}
			bl = &pem.Block{Type: tp, Bytes: by}
		}
		if err := pem.Encode(pe.out, bl); err != nil {
			return err
		}
	}
	return nil
}

// PEMMarshaler marshals itself into a PEM block
type PEMMarshaler interface {
	MarshalPEM() (*pem.Block, error)
}

// PEMUnmarshaler unmarshals a PEM block into itself
type PEMUnmarshaler interface {
	UnmarshalPEM(bl *pem.Block) error
}
