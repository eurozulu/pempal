package fileformats

import (
	"encoding/pem"
	"fmt"
	"gopkg.in/yaml.v3"
)

const PEM_TEMPLATE = "TEMPLATE"

type yamlTemplateReader struct {
}

func (y yamlTemplateReader) Unmarshal(by []byte) ([]*pem.Block, error) {
	var pt pemOnly
	if err := yaml.Unmarshal(by, &pt); err != nil {
		return nil, fmt.Errorf("invalid template, %w", err)
	}
	if pt.PemType == "" {
		return nil, fmt.Errorf("invalid template, missing 'pem_type'")
	}
	return []*pem.Block{{
		Type:  fmt.Sprintf("%s %s", pt.PemType, PEM_TEMPLATE),
		Bytes: by,
	}}, nil
}

type pemOnly struct {
	PemType string `yaml:"pem_type,omitempty"`
}
