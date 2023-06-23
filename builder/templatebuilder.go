package builder

import (
	"bytes"
	"github.com/eurozulu/pempal/model"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"gopkg.in/yaml.v2"
)

type TemplateBuilder []templates.Template

func (tb TemplateBuilder) MergeTemplates() (templates.Template, error) {
	m := utils.FlatMap{}
	if err := tb.Apply(&m); err != nil {
		return nil, err
	}
	return TemplateFromValue(&m)
}

func (tb TemplateBuilder) Apply(v interface{}) error {
	for _, t := range tb {
		if err := t.Apply(v); err != nil {
			return err
		}
	}
	return nil
}

func (tb TemplateBuilder) ResourceType() model.ResourceType {
	for _, t := range tb {
		s := equalsResourceTemplate([]byte(t.String()))
		if s == "" {
			continue
		}
		return model.ParseResourceType(s)
	}
	return model.Unknown
}

func equalsResourceTemplate(data []byte) string {
	for k, v := range model.DefaultResourceTemplates {
		if bytes.Equal(data, v) {
			return k
		}
	}
	return ""
}

func TemplateFromValue(v interface{}) (templates.Template, error) {
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return templates.NewTemplate(buf.Bytes())
}
