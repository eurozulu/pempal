package commands

import (
	"github.com/go-yaml/yaml"
	"io"
	"pempal/resourceio"
)

type TemplatesCommand struct {
	Name         string `flag:"name"`
	TemplatePath string `flag:"template-path,path"`
}

func (t TemplatesCommand) Execute(args []string, out io.Writer) error {
	if t.TemplatePath == "" {
		t.TemplatePath = configuration.TemplatePath
	}

	m := map[string]interface{}{}
	if err := resourceio.MergeTemplatesInto(&m, t.TemplatePath, args...); err != nil {
		return err
	}
	by, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	by = append(by, '\n')
	_, err = out.Write(by)
	return err
}
