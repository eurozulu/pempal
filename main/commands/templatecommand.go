package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"os"
	"pempal/resourceio"
	"pempal/templates"
)

var templateManager templates.TemplateManager

type TemplateCommand struct {
	TemplatePath string `flag:"template-path,path"`
	Raw          bool   `flag:"raw"`
	Add          string `flag:"add"`
	Remove       string `flag:"remove"`
}

func (cmd TemplateCommand) Execute(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide one or more template names")
	}
	if err := cmd.setTemplateManager(); err != nil {
		return err
	}

	if cmd.Remove != "" {
		return cmd.removeTemplates(args)
	}
	var data []byte
	var err error
	if !cmd.Raw {
		data, err = mergeNamedTemplates(args)
	} else {
		data, err = rawNamedTemplates(args)
	}
	if cmd.Add != "" {
		return cmd.addTemplate(cmd.Add, data)
	}

	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = out.Write(data)
	return err
}

func (cmd TemplateCommand) addTemplate(name string, data []byte) error {
	t, err := templateManager.ParseTemplate(data)
	if err != nil {
		return err
	}
	return templateManager.AddTemplate(name, t)
}

func (cmd TemplateCommand) removeTemplates(names []string) error {
	for _, name := range names {
		if err := templateManager.RemoveTemplate(name); err != nil {
			return err
		}
	}
	return nil
}

func (cmd TemplateCommand) setTemplateManager() error {
	if cmd.TemplatePath == "" {
		cmd.TemplatePath = configuration.TemplatePath
	}
	tm, err := resourceio.NewResourceTemplateManager(cmd.TemplatePath)
	if err != nil {
		return err
	}
	templateManager = tm
	return nil
}

func mergeNamedTemplates(names []string) ([]byte, error) {
	m := map[string]interface{}{}
	if err := templateManager.MergeTemplatesInto(&m, names...); err != nil {
		return nil, err
	}
	return yaml.Marshal(m)
}

func rawNamedTemplates(names []string) ([]byte, error) {
	temps, err := templateManager.TemplatesByName(names...)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	for i, t := range temps {
		if i > 0 {
			buf.WriteRune('\n')
		}
		buf.Write(t.Raw())
	}
	return buf.Bytes(), nil
}

func readStdIn() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		buf.Write(scan.Bytes())
	}
	if scan.Err() != nil && scan.Err() != io.EOF {
		return nil, scan.Err()
	}
	return buf.Bytes(), nil
}
