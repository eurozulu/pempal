package commands

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"pempal/resourceio"
)

type TemplateCommand struct {
	TemplatePath string `flag:"template-path,path"`
}

func (t TemplateCommand) Execute(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide one or more template names")
	}
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
