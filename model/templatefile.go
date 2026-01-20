package model

import (
	"bytes"
	"github.com/eurozulu/pempal/tools"
	"gopkg.in/yaml.v2"
	"strings"
)

const ExtendsComment = "extends"

// TemplateFile represents a template from a file. i.e. NOT a Base template.
// Template files may contain an #extends to name additional templates.
// Templates are formatted as yaml
type TemplateFile struct {
	Path string
	Data []byte
}

func (t TemplateFile) Extends() []string {
	ext, _ := splitExtends(t.Data)
	return ext
}

func (t *TemplateFile) MarshalBinary() (data []byte, err error) {
	return t.Data, nil
}

func (t *TemplateFile) UnmarshalBinary(data []byte) error {
	t.Data = data
	return nil
}

func (t TemplateFile) String() string {
	return string(t.Data)
}

func (t TemplateFile) CleanData() []byte {
	_, data := splitExtends(t.Data)
	return data
}

func (t TemplateFile) IsValid() error {
	data := t.CleanData()
	m := map[string]interface{}{}
	return yaml.Unmarshal(data, &m)
}

func splitExtends(data []byte) ([]string, []byte) {
	if !bytes.HasPrefix(data, []byte("#")) {
		return nil, data
	}
	b := bytes.TrimSpace(data[1:])
	if !bytes.HasPrefix(b, []byte(ExtendsComment)) {
		return nil, data
	}
	b = bytes.TrimSpace(bytes.TrimPrefix(b, []byte(ExtendsComment)))
	bb := bytes.SplitN(b, []byte("\n"), 2)
	names := tools.TrimSlice(strings.Split(string(bb[0]), ","))
	var dat []byte
	if len(bb) > 1 {
		dat = bytes.TrimSpace(bb[1])
	}
	return names, dat
}
