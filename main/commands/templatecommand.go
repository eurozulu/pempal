package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"os"
)

type TemplateCommand struct {
	Raw    bool   `flag:"raw"`
	Add    string `flag:"add"`
	Remove string `flag:"remove"`
}

func (cmd TemplateCommand) Execute(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide one or more template names")
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
	t, err := ResourceTemplates.ParseTemplate(data)
	if err != nil {
		return err
	}
	return ResourceTemplates.AddTemplate(name, t)
}

func (cmd TemplateCommand) removeTemplates(names []string) error {
	for _, name := range names {
		if err := ResourceTemplates.RemoveTemplate(name); err != nil {
			return err
		}
	}
	return nil
}

func mergeNamedTemplates(names []string) ([]byte, error) {
	m := map[string]interface{}{}
	if err := ResourceTemplates.MergeTemplatesInto(&m, names...); err != nil {
		return nil, err
	}
	return yaml.Marshal(m)
}

func rawNamedTemplates(names []string) ([]byte, error) {
	temps, err := ResourceTemplates.TemplatesByName(names...)
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
