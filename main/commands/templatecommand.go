package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"pempal/templates"
	"strings"
)

const (
	formatYAML = "yaml"
	formatJSON = "json"
	formatRAW  = "raw"
)

var formats = []string{formatYAML, formatJSON, formatRAW}

type TemplateCommand struct {
	Format string `flag:"format"`
}

func (cmd TemplateCommand) Execute(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("must provide one or more template names")
	}

	temps, err := parseArgumentsToTemplates(args)
	if err != nil {
		return err
	}

	switch cmd.resolveFormat() {
	case formatYAML:
		return outputAsYaml(temps, out)
	case formatJSON:
		return outputAsJson(temps, out)
	case formatRAW:
		return outputAsRaw(temps, out)
	default:
		return fmt.Errorf("%s is an unknown format. Use one of %v", cmd.Format, formats)
	}
}

func (cmd TemplateCommand) resolveFormat() string {
	if cmd.Format == "" {
		return formats[0]
	}
	for _, f := range formats {
		if strings.EqualFold(f, cmd.Format) {
			return f
		}
	}
	return cmd.Format
}

func outputAsYaml(temps []templates.Template, out io.Writer) error {
	m := map[string]interface{}{}
	if err := templates.ApplyTemplatesTo(&m, temps); err != nil {
		return err
	}
	return yaml.NewEncoder(out).Encode(m)
}

func outputAsJson(temps []templates.Template, out io.Writer) error {
	m := map[string]interface{}{}
	if err := templates.ApplyTemplatesTo(&m, temps); err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(m)
}

func outputAsRaw(temps []templates.Template, out io.Writer) error {
	buf := bytes.NewBuffer(nil)
	for i, t := range temps {
		if i > 0 {
			buf.WriteRune('\n')
		}
		buf.Write(t.Raw())
	}
	_, err := out.Write(buf.Bytes())
	return err
}

func parseArgumentsToTemplates(args []string) ([]templates.Template, error) {
	if ResourceTemplates == nil {
		return nil, fmt.Errorf("template manager unavailable.")
	}
	var temps []templates.Template
	for _, arg := range args {
		t, err := ResourceTemplates.TemplatesByName(arg)
		if err != nil {
			return nil, err
		}
		temps = append(temps, t...)
	}
	return temps, nil
}
