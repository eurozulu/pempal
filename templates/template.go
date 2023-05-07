package templates

import (
	"fmt"
	"github.com/go-yaml/yaml"
)

// Template represents a raw byte slice of data in a specific format.
// Template may optionally be preceeded with #tags, named tags specifying
// which other resources are related to this template.
type Template interface {
	// Raw returns the original, raw bytes of the template, prior to any parsing and attribution
	Raw() []byte

	// String returns a string of the template bytes after attribution (tags applied)
	fmt.Stringer

	// Tags returns any #tags found in this template
	Tags() Tags

	// Apply this template to the given object.
	// This unmarshalls the template int the given object.
	// Object can be a map with string/interface key/values or
	// a struct with public fields matching the properties in the template
	Apply(in interface{}) error
}

type yamlTemplate struct {
	raw    []byte
	tags   Tags
	parsed []byte
}

func (t yamlTemplate) Raw() []byte {
	return t.raw
}

func (t yamlTemplate) String() string {
	return string(t.parsed)
}

func (t yamlTemplate) Tags() Tags {
	return t.tags
}

// Apply applies this template to the given object
func (t yamlTemplate) Apply(out interface{}) error {
	return yaml.Unmarshal(t.parsed, out)
}

func (t *yamlTemplate) extendTemplates(temps []Template) error {
	m := map[string]interface{}{}
	if err := ApplyTemplatesTo(&m, temps); err != nil {
		return err
	}
	if err := t.Apply(&m); err != nil {
		return err
	}
	by, err := yaml.Marshal(&m)
	if err != nil {
		return err
	}
	t.parsed = by
	return nil
}

func newYamlTemplate(raw []byte, tags Tags, parsed []byte, extends []Template) (Template, error) {
	t := &yamlTemplate{
		raw:    raw,
		parsed: parsed,
		tags:   tags,
	}
	if err := t.extendTemplates(extends); err != nil {
		return nil, err
	}
	return t, nil
}
