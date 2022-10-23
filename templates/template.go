package templates

import (
	"bytes"
	"crypto/x509"
	"gopkg.in/yaml.v3"
)

// Template is a set of certificate properties which can be applied to a certiifcate.
type Template interface {
	Apply(c *x509.Certificate) error
}

type yamlTemplate struct {
	raw []byte
}

func (t yamlTemplate) Apply(c *x509.Certificate) error {
	return yaml.NewDecoder(bytes.NewBuffer(t.raw)).Decode(c)
}

func (t yamlTemplate) MarshalYAML() (interface{}, error) {
	return t.raw, nil
}

func NewTemplate(by []byte) Template {
	return &yamlTemplate{raw: by}
}
