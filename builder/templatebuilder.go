package builder

import (
	"fmt"
	"github.com/eurozulu/pempal/logger"
	"github.com/eurozulu/pempal/templates"
	"github.com/eurozulu/pempal/utils"
	"github.com/go-yaml/yaml"
)

func MergeTemplates(temps []templates.Template) (templates.Template, error) {
	m := utils.FlatMap{}
	if err := ApplyTemplates(&m, temps); err != nil {
		return nil, err
	}
	data, err := yaml.Marshal(&m)
	if err != nil {
		return nil, err
	}
	return templates.NewTemplate(data)
}

func ApplyTemplates(v interface{}, temps []templates.Template) error {
	for _, t := range temps {
		var err error
		data := t.Bytes()
		if containsGoTemplates(data) {
			logger.Debug("go template detected.  executing template engine")
			data, err = executeGoTemplate(data, v)
			if err != nil {
				return fmt.Errorf("failed to execute template %v", err)
			}
		}
		if err = yaml.Unmarshal(data, v); err != nil {
			return err
		}
	}
	return nil
}
