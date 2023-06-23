package templates

import (
	"bytes"
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/go-yaml/yaml"
)

// Template represents a collection of named properties.
type Template interface {

	// Tags returns any #tags found in this template
	Tags() Tags

	// Bytes returns the 'raw', unformatted template, including tags
	Bytes() []byte

	Apply(v interface{}) error

	// String returns a string of the formatted template after tags and macros applied
	fmt.Stringer
}

type yamlTemplate struct {
	raw  []byte
	tags Tags
}

func (t yamlTemplate) Tags() Tags {
	return t.tags
}

func (t yamlTemplate) Apply(v interface{}) error {
	var err error
	data := t.Bytes()
	if containsGoTemplates(data) {
		logger.Debug("go template detected.  executing template engine")
		data, err = executeGoTemplate(data, v)
		if err != nil {
			return fmt.Errorf("failed to execute template %v", err)
		}
	}
	return yaml.Unmarshal(data, v)
}

func (t yamlTemplate) String() string {
	buf := bytes.NewBuffer(nil)
	for _, tg := range t.tags {
		buf.WriteString(tg.String())
		buf.WriteRune('\n')
	}
	buf.Write(t.raw)
	return buf.String()
}

func (t yamlTemplate) Bytes() []byte {
	return t.raw
}

func NewTemplate(data []byte) (Template, error) {
	tags, raw, err := parseTags(data)
	if err != nil {
		return nil, err
	}
	return &yamlTemplate{
		tags: tags,
		raw:  raw,
	}, nil
}
