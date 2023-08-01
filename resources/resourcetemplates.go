package resources

import (
	"bytes"
	"github.com/eurozulu/pempal/templates"
	"github.com/go-yaml/yaml"
	"strings"
)

func NewResourceTemplatesManager(storePath []string) (templates.TemplateManager, error) {
	tm, err := templates.NewTemplateManager(storePath)
	if err != nil {
		return nil, err
	}
	if err = addResourceTemplates(tm); err != nil {
		return nil, err
	}
	return tm, nil
}

func resourceTemplateAsBytes(rt ResourceType) ([]byte, error) {
	t, err := ResourceTemplateByType(rt)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	if err = yaml.NewEncoder(buf).Encode(&t); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func addResourceTemplates(tm templates.TemplateManager) error {
	for _, rt := range resourceTypes[1:] {
		data, err := resourceTemplateAsBytes(rt)
		if err != nil {
			return err
		}
		if err := tm.AddDefaultTemplate(strings.ToLower(rt.String()), data); err != nil {
			return err
		}
	}
	return nil
}
