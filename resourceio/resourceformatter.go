package resourceio

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"github.com/eurozulu/pempal/model"
	"github.com/go-yaml/yaml"
)

// ResourceFormatter formats Resources into a specific format.
// Available formats are PEM, DER or YAML
type ResourceFormatter interface {
	FormatResources(res ...model.Resource) ([]byte, error)
}

// PemFormatter formats bytes into one or more Pem Blocks
type PemFormatter struct{}

func (p PemFormatter) FormatResources(res ...model.Resource) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, r := range res {
		if _, err := buf.WriteString(r.String()); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// DerFormatter will transform a single Resource into its DER encoded bytes
type DerFormatter struct{}

func (p DerFormatter) FormatResources(res ...model.Resource) ([]byte, error) {
	if len(res) > 1 {
		return nil, fmt.Errorf("FormatDER can only format single resources.")
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("nothing to format")
	}
	blk, _ := pem.Decode([]byte(res[0].String()))
	if blk == nil {
		return nil, fmt.Errorf("failed to parse the pem of resource")
	}
	return blk.Bytes, nil
}

// YamlResourceFormatter formats the resources into a yaml document.
// Multiple resources are delimited with the yaml page break '---'
// Transformation takes place via the appropriate ResourceDTO
type YamlResourceFormatter struct{}

func (t YamlResourceFormatter) FormatResources(res ...model.Resource) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for i, r := range res {
		dto, err := model.DTOForResource(r)
		if err != nil {
			return nil, err
		}
		data, err := yaml.Marshal(dto)
		if err != nil {
			return nil, err
		}

		if i > 0 {
			// write yaml document seperator
			if _, err = buf.WriteString("---\n"); err != nil {
				return nil, err
			}
		}
		if _, err = buf.Write(data); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// NewResourceFormatter creates a ResourceFormatter for the given type.
// type MUST be one of: FormatPEM, FormatDER, FormatYAML
// Anything else will cause a panic.
func NewResourceFormatter(format ResourceFormat) ResourceFormatter {
	switch format {
	case FormatPEM:
		return PemFormatter{}
	case FormatDER:
		return DerFormatter{}
	case FormatYAML:
		return YamlResourceFormatter{}
	default:
		panic("unknown resource format")
	}
}
